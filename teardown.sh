#!/bin/bash

export MSYS_NO_PATHCONV=1
starttime=$(date +%s)

pushd ./basic-network
./teardown.sh
popd

cat <<EOF

Total setup execution time : $(($(date +%s) - starttime)) secs ...

EOF