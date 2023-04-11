package test

import (
	"testing"

	"github.com/sllt/dtm/client/dtmcli"
	"github.com/sllt/dtm/client/dtmcli/dtmimp"
	"github.com/sllt/dtm/client/dtmgrpc"
	"github.com/sllt/dtm/dtmsvr/storage/sql"
	"github.com/sllt/dtm/dtmutil"
	"github.com/stretchr/testify/assert"
)

func TestGeneralDB(t *testing.T) {
	if conf.Store.IsDB() {
		testSql(t)
		testDbAlone(t)
	}
}

func testSql(t *testing.T) {
	conf := conf.Store.GetDBConf()
	conf.Host = "127.0.0.1" // use a new host to trigger SetDBConn called
	db := dtmutil.DbGet(conf, sql.SetDBConn)
	err := func() (rerr error) {
		defer dtmimp.P2E(&rerr)
		db.Must().Exec("select a")
		return nil
	}()
	assert.NotEqual(t, nil, err)
}

func testDbAlone(t *testing.T) {
	db, err := dtmimp.StandaloneDB(conf.Store.GetDBConf())
	assert.Nil(t, err)
	_, err = dtmimp.DBExec(conf.Store.Driver, db, "select 1")
	assert.Equal(t, nil, err)
	_, err = dtmimp.DBExec(conf.Store.Driver, db, "")
	assert.Equal(t, nil, err)
	db.Close()
	_, err = dtmimp.DBExec(conf.Store.Driver, db, "select 1")
	assert.NotEqual(t, nil, err)
}

func TestMustGenGid(t *testing.T) {
	dtmgrpc.MustGenGid(dtmutil.DefaultGrpcServer)
	dtmcli.MustGenGid(dtmutil.DefaultHTTPServer)
}
