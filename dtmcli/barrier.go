/*
 * Copyright (c) 2021 yedf. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package dtmcli

import (
	"database/sql"
	"fmt"
	"net/url"

	"github.com/dtm-labs/dtm/dtmcli/dtmimp"
	"github.com/dtm-labs/dtm/dtmcli/logger"
)

// BarrierBusiFunc type for busi func
type BarrierBusiFunc func(tx *sql.Tx) error

// BranchBarrier every branch info
type BranchBarrier struct {
	TransType string
	Gid       string
	BranchID  string
	Op        string
	BarrierID int
}

func (bb *BranchBarrier) String() string {
	return fmt.Sprintf("transInfo: %s %s %s %s", bb.TransType, bb.Gid, bb.BranchID, bb.Op)
}

func (bb *BranchBarrier) newBarrierID() string {
	bb.BarrierID++
	return fmt.Sprintf("%02d", bb.BarrierID)
}

// BarrierFromQuery construct transaction info from request
func BarrierFromQuery(qs url.Values) (*BranchBarrier, error) {
	return BarrierFrom(qs.Get("trans_type"), qs.Get("gid"), qs.Get("branch_id"), qs.Get("op"))
}

// BarrierFrom construct transaction info from request
func BarrierFrom(transType, gid, branchID, op string) (*BranchBarrier, error) {
	ti := &BranchBarrier{
		TransType: transType,
		Gid:       gid,
		BranchID:  branchID,
		Op:        op,
	}
	if ti.TransType == "" || ti.Gid == "" || ti.BranchID == "" || ti.Op == "" {
		return nil, fmt.Errorf("invalid trans info: %v", ti)
	}
	return ti, nil
}

func insertBarrier(tx DB, transType string, gid string, branchID string, op string, barrierID string, reason string) (int64, error) {
	if op == "" {
		return 0, nil
	}
	sql := dtmimp.GetDBSpecial().GetInsertIgnoreTemplate(dtmimp.BarrierTableName+"(trans_type, gid, branch_id, op, barrier_id, reason) values(?,?,?,?,?,?)", "uniq_barrier")
	return dtmimp.DBExec(tx, sql, transType, gid, branchID, op, barrierID, reason)
}

// Call 子事务屏障，详细介绍见 https://zhuanlan.zhihu.com/p/388444465
// tx: 本地数据库的事务对象，允许子事务屏障进行事务操作
// busiCall: 业务函数，仅在必要时被调用
func (bb *BranchBarrier) Call(tx *sql.Tx, busiCall BarrierBusiFunc) (rerr error) {
	bid := bb.newBarrierID()
	defer dtmimp.DeferDo(&rerr, func() error {
		return tx.Commit()
	}, func() error {
		return tx.Rollback()
	})
	originOp := map[string]string{
		BranchCancel:     BranchTry,
		BranchCompensate: BranchAction,
	}[bb.Op]

	originAffected, oerr := insertBarrier(tx, bb.TransType, bb.Gid, bb.BranchID, originOp, bid, bb.Op)
	currentAffected, rerr := insertBarrier(tx, bb.TransType, bb.Gid, bb.BranchID, bb.Op, bid, bb.Op)
	logger.Debugf("originAffected: %d currentAffected: %d", originAffected, currentAffected)

	if rerr == nil && bb.Op == "msg" && currentAffected == 0 { // for msg's DoAndSubmit, repeated insert should be rejected.
		return ErrDuplicated
	}

	if rerr == nil {
		rerr = oerr
	}

	if (bb.Op == BranchCancel || bb.Op == BranchCompensate) && originAffected > 0 || // null compensate
		currentAffected == 0 { // repeated request or dangled request
		return
	}
	if rerr == nil {
		rerr = busiCall(tx)
	}
	return
}

// CallWithDB the same as Call, but with *sql.DB
func (bb *BranchBarrier) CallWithDB(db *sql.DB, busiCall BarrierBusiFunc) error {
	tx, err := db.Begin()
	if err == nil {
		err = bb.Call(tx, busiCall)
	}
	return err
}

// QueryPrepared queries prepared data
func (bb *BranchBarrier) QueryPrepared(db *sql.DB) error {
	_, err := insertBarrier(db, bb.TransType, bb.Gid, "00", "msg", "01", "rollback")
	var reason string
	if err == nil {
		sql := fmt.Sprintf("select reason from %s where gid=? and branch_id=? and op=? and barrier_id=?", dtmimp.BarrierTableName)
		err = db.QueryRow(sql, bb.Gid, "00", "msg", "01").Scan(&reason)
	}
	if reason == "rollback" {
		return ErrFailure
	}
	return err
}
