package test

import (
	"testing"

	"github.com/sllt/dtm/client/dtmcli"
	"github.com/sllt/dtm/client/dtmcli/dtmimp"
	"github.com/sllt/dtm/dtmsvr"
	"github.com/sllt/dtm/dtmutil"
	"github.com/sllt/dtm/test/busi"
	"github.com/stretchr/testify/assert"
)

func genMsgDelay(gid string) *dtmcli.Msg {
	req := busi.GenReqHTTP(30, false, false)
	msg := dtmcli.NewMsg(dtmutil.DefaultHTTPServer, gid).
		Add(busi.Busi+"/TransOut", &req).
		Add(busi.Busi+"/TransIn", &req).SetDelay(10)
	msg.QueryPrepared = busi.Busi + "/QueryPrepared"
	return msg
}

func TestMsgDelayNormal(t *testing.T) {
	gid := dtmimp.GetFuncName()
	msg := genMsgDelay(gid)
	submitForwardCron(0, func() {
		msg.Submit()
		waitTransProcessed(msg.Gid)
	})

	dtmsvr.NowForwardDuration = 0
	assert.Equal(t, []string{StatusPrepared, StatusPrepared}, getBranchesStatus(msg.Gid))
	assert.Equal(t, StatusSubmitted, getTransStatus(msg.Gid))
	cronTransOnceForwardCron(t, "", 0)
	cronTransOnceForwardCron(t, "", 8)
	cronTransOnceForwardCron(t, gid, 12)
	assert.Equal(t, []string{StatusSucceed, StatusSucceed}, getBranchesStatus(msg.Gid))
	assert.Equal(t, StatusSucceed, getTransStatus(msg.Gid))
}
