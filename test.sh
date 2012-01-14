#!/bin/bash

source test/helper/http.sh

./_bin/webreduced -listen $addr &
server_pid=$!

# XXX find better way act on server start
sleep 1

roundup test/system/apps-test.sh

kill $server_pid
