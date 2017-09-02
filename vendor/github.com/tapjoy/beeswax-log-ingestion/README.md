Production URL: https://beeswax.tapjoy.com/wins

Dependencies:
`sudo pip install requests`

## Environment

Setup a `.env` in your local environment, which Docker will read on turnup. See `.env-example` for an example `.env`. One of the uses of this would be for setting up `NEW_RELIC_APP_NAME` if you want to report to New Relic while developing.

## Run python script on _local_ machine that generates protobuf requests to be sent to dev app

Dependencies: `sudo pip install requests`

```bash
./beeswax/tools/win_events/win_log_requester/win_events/win_log_requester beeswax/tools/win_events/win_log_requester/sample_ad_log.txt http://localhost:32785/wins --path-to-responses-file beeswax/tools/win_events/win_log_requester/ad_log_response.txt --log-level debug
```

test requests from beeswax:
`./beeswax/tools/win_events/win_log_requester/win_events/win_log_requester beeswax/tools/win_events/win_log_requester/sample_ad_log.txt http://localhost:32785/wins --path-to-responses-file beeswax/tools/win_events/win_log_requester/ad_log_response.txt --log-level debug`

## generate protobuf files

if you are missing a dependency: run `gvt restore`

proto gen:

rm -rf $GOPATH/src/beeswax/ && mkdir $GOPATH/src/beeswax && cp -r $GOPATH/src/github.com/tapjoy/beeswax-log-ingestion/beeswax/* $GOPATH/src/beeswax && protoc --go_out=$GOPATH/src/github.com/tapjoy/beeswax-log-ingestion/protos beeswax/*/*.proto

add a dependency:

govendor fetch github.com/golang/protobuf/proto
