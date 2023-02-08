#!/bin/bash

HOST=localhost
PORT=8080
TEST_TIME=60s
SLEEP_TIME=10
ITERATION=5
CONCURENCY=(200 400 600 800 1000)
AUTH_METHOD=basic # basic / jwt
USERNAME=perftest
PASSWORD=perftest

init_db() {
   PGPASSWORD=postgres psql -h $HOST -U postgres postgres < dump.sql
}

drop_table_db() {
  PGPASSWORD=postgres psql -h $HOST -U postgres postgres -c "DROP TABLE user_person, hardware, node, sensor, channel;"
}

rollback_db() {
  drop_table_db
  init_db
}

perftest() {
  local method=$1
  local endpoint=$2
  local data=$3
  local do_rollback=$4
  for concurrent in "${CONCURENCY[@]}"
  do
    for i in $(seq 1 $ITERATION)
    do
      echo "[$method $endpoint $concurrent][$i]"
      if [ $method == "GET" ]
      then
        hey -n $concurrent -z $TEST_TIME -H "$HEADER" http://$HOST:$PORT$endpoint
      elif [ $method == "PUT" ] || [ $method == "POST" ]
      then
        hey -m $method -T "appliaction/json" -d "$data" -n $concurrent -z $TEST_TIME -H "$HEADER" http://$HOST:$PORT$endpoint
      fi
      sleep $SLEEP_TIME
      if $do_rollback
      then
        rollback_db > /dev/null
      fi
    done
  done
}

auth() {
  if [ $AUTH_METHOD == "jwt" ]
  then
    HEADER="Authorization: Bearer $(http $HOST:$PORT/user/login username=perftest password=perftest | jq -r '.token')"
  elif [ $AUTH_METHOD == "basic" ]
  then
    HEADER="Authorization: Basic $(echo -n $USERNAME:$PASSWORD | base64)"
  fi
}

rollback_db > /dev/null
auth

## perftest METHOD ENDPOINT DATA DO_ROLLBACK
perftest "GET" "/node" "" false
# perftest "POST" "/channel" "{\"value\": 1.33, \"id_sensor\": 1}" true
# perftest "GET" "/node/1" ""  false
# perftest "PUT" "/node/1" "{ \"name\":\"test\",\"location\":\"test\",\"id_hardware\":1 }" true
