#!/bin/bash
set -e

export MSYS_NO_PATHCONV=1
starttime=$(date +%s)

docker-compose -f docker-compose-simple.yaml down

docker-compose -f docker-compose-simple.yaml up -d

sleep 10

docker exec chaincode sh -c "./build.sh"

docker exec chaincode sh -c "./start.sh"

cat <<EOF

Total setup execution time : $(($(date +%s) - starttime)) secs ...

EOF