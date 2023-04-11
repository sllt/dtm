package test

import (
	"testing"

	"github.com/sllt/dtm/client/dtmcli/dtmimp"
	"github.com/sllt/dtm/test/busi"
	"github.com/stretchr/testify/assert"
)

func TestMsgWebhook(t *testing.T) {
	msg := genMsg(dtmimp.GetFuncName())
	busi.MainSwitch.TransInResult.SetOnce("ERROR")
	msg.Submit()
	assert.Equal(t, StatusSubmitted, getTransStatus(msg.Gid))
	waitTransProcessed(msg.Gid)
	busi.MainSwitch.TransInResult.SetOnce("ERROR")
	cronTransOnce(t, msg.Gid)
	busi.MainSwitch.TransInResult.SetOnce("ERROR")
	cronTransOnce(t, msg.Gid)
	assert.Equal(t, msg.Gid, busi.WebHookResult["gid"])
	cronTransOnce(t, msg.Gid)
	assert.Equal(t, []string{StatusSucceed, StatusSucceed}, getBranchesStatus(msg.Gid))
	assert.Equal(t, StatusSucceed, getTransStatus(msg.Gid))
}

func TestMsgWebhookError(t *testing.T) {
	msg := genMsg(dtmimp.GetFuncName())
	busi.MainSwitch.TransInResult.SetOnce("ERROR")
	msg.Submit()
	assert.Equal(t, StatusSubmitted, getTransStatus(msg.Gid))
	waitTransProcessed(msg.Gid)
	busi.MainSwitch.TransInResult.SetOnce("ERROR")
	cronTransOnce(t, msg.Gid)
	busi.MainSwitch.TransInResult.SetOnce("ERROR")
	cronTransOnce(t, msg.Gid)
	assert.Equal(t, msg.Gid, busi.WebHookResult["gid"])
	cronTransOnce(t, msg.Gid)
	assert.Equal(t, []string{StatusSucceed, StatusSucceed}, getBranchesStatus(msg.Gid))
	assert.Equal(t, StatusSucceed, getTransStatus(msg.Gid))
}
