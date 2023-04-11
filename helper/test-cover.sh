set -x
export DTM_DEBUG=1
echo "mode: count" > coverage.txt
for store in redis boltdb mysql postgres; do
  TEST_STORE=$store go test -failfast -covermode count -coverprofile=profile.out -coverpkg=github.com/sllt/dtm/client/dtmcli,github.com/sllt/dtm/client/dtmcli/dtmimp,github.com/dtm-labs/logger,github.com/sllt/dtm/client/dtmgrpc,github.com/sllt/dtm/client/workflow,github.com/sllt/dtm/client/dtmgrpc/dtmgimp,github.com/sllt/dtm/dtmsvr,dtmsvr/config,github.com/sllt/dtm/dtmsvr/storage,github.com/sllt/dtm/dtmsvr/storage/boltdb,github.com/sllt/dtm/dtmsvr/storage/redis,github.com/sllt/dtm/dtmsvr/storage/registry,github.com/sllt/dtm/dtmsvr/storage/sql,github.com/sllt/dtm/dtmutil -gcflags=-l ./... || exit 1
    echo "TEST_STORE=$store finished"
    if [ -f profile.out ]; then
        cat profile.out | grep -v 'mode:' >> coverage.txt
        echo > profile.out
    fi
done

# go tool cover -html=coverage.txt

# curl -s https://codecov.io/bash | bash
