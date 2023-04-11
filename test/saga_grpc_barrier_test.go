/*
 * Copyright (c) 2021 yedf. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package test

import (
	"testing"

	"github.com/sllt/dtm/client/dtmcli/dtmimp"
	"github.com/sllt/dtm/client/dtmgrpc"
	"github.com/sllt/dtm/dtmutil"
	"github.com/sllt/dtm/test/busi"
	"github.com/stretchr/testify/assert"
)

func TestSagaGrpcBarrierNormal(t *testing.T) {
	saga := genSagaGrpcBarrier(dtmimp.GetFuncName(), false, false)
	err := saga.Submit()
	assert.Nil(t, err)
	waitTransProcessed(saga.Gid)
	assert.Equal(t, StatusSucceed, getTransStatus(saga.Gid))
	assert.Equal(t, []string{StatusPrepared, StatusSucceed, StatusPrepared, StatusSucceed}, getBranchesStatus(saga.Gid))
}

func TestSagaGrpcBarrierRollback(t *testing.T) {
	saga := genSagaGrpcBarrier(dtmimp.GetFuncName(), false, true)
	err := saga.Submit()
	assert.Nil(t, err)
	waitTransProcessed(saga.Gid)
	assert.Equal(t, StatusFailed, getTransStatus(saga.Gid))
	assert.Equal(t, []string{StatusSucceed, StatusSucceed, StatusSucceed, StatusFailed}, getBranchesStatus(saga.Gid))
}

func genSagaGrpcBarrier(gid string, outFailed bool, inFailed bool) *dtmgrpc.SagaGrpc {
	saga := dtmgrpc.NewSagaGrpc(dtmutil.DefaultGrpcServer, gid)
	req := busi.GenReqGrpc(30, outFailed, inFailed)
	saga.Add(busi.BusiGrpc+"/busi.Busi/TransOutBSaga", busi.BusiGrpc+"/busi.Busi/TransOutRevertBSaga", req)
	saga.Add(busi.BusiGrpc+"/busi.Busi/TransInBSaga", busi.BusiGrpc+"/busi.Busi/TransInRevertBSaga", req)
	return saga
}
