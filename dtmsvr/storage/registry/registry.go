package registry

import (
	"time"

	"github.com/dtm-labs/logger"

	"github.com/sllt/dtm/dtmsvr/config"
	"github.com/sllt/dtm/dtmsvr/storage"
	"github.com/sllt/dtm/dtmsvr/storage/boltdb"
	"github.com/sllt/dtm/dtmsvr/storage/redis"
	"github.com/sllt/dtm/dtmsvr/storage/sql"
)

var conf = &config.Config

// StorageFactory is factory to get storage instance.
type StorageFactory interface {
	// GetStorage will return the Storage instance.
	GetStorage() storage.Store
}

var sqlFac = &SingletonFactory{
	creatorFunction: func() storage.Store {
		return &sql.Store{}
	},
}

var storeFactorys = map[string]StorageFactory{
	"boltdb": &SingletonFactory{
		creatorFunction: func() storage.Store {
			return boltdb.NewStore(conf.Store.DataExpire, conf.RetryInterval)
		},
	},
	"redis": &SingletonFactory{
		creatorFunction: func() storage.Store {
			return &redis.Store{}
		},
	},
	"mysql":    sqlFac,
	"postgres": sqlFac,
}

// GetStore returns storage.Store
func GetStore() storage.Store {
	return storeFactorys[conf.Store.Driver].GetStorage()
}

// WaitStoreUp wait for db to go up
func WaitStoreUp() {
	for err := GetStore().Ping(); err != nil; err = GetStore().Ping() {
		logger.Infof("wait store up: %v", err)
		time.Sleep(3 * time.Second)
	}
}
