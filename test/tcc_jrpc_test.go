package test

import (
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/sllt/dtm/client/dtmcli"
	"github.com/sllt/dtm/client/dtmcli/dtmimp"
	"github.com/sllt/dtm/dtmutil"
	"github.com/sllt/dtm/test/busi"
	"github.com/stretchr/testify/assert"
)

func TestTccJrpcNormal(t *testing.T) {
	req := busi.GenReqHTTP(30, false, false)
	gid := dtmimp.GetFuncName()
	err := dtmcli.TccGlobalTransaction2(dtmutil.DefaultJrpcServer, gid, func(tcc *dtmcli.Tcc) {
		tcc.Protocol = dtmimp.Jrpc
	}, func(tcc *dtmcli.Tcc) (*resty.Response, error) {
		_, err := tcc.CallBranch(req, Busi+"/TransOut", Busi+"/TransOutConfirm", Busi+"/TransOutRevert")
		assert.Nil(t, err)
		return tcc.CallBranch(req, Busi+"/TransIn", Busi+"/TransInConfirm", Busi+"/TransInRevert")
	})
	assert.Nil(t, err)
	waitTransProcessed(gid)
	assert.Equal(t, StatusSucceed, getTransStatus(gid))
	assert.Equal(t, []string{StatusPrepared, StatusSucceed, StatusPrepared, StatusSucceed}, getBranchesStatus(gid))
}
