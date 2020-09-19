#!/bin/bash

set -e

if [ $# -ne 1 ]; then
    echo "$0 <count>"
    exit 1
fi

index=0
while [ $index -lt $1 ]; do
    lb=`uuidgen | tr '[:upper:]' '[:lower:]'`
    bp=`uuidgen | tr '[:upper:]' '[:lower:]'`
    curl -s -X POST "http://localhost:8080/mapping" -d \
    "{ \
        \"loadbalancer\": \"$lb\", \
        \"bigip\":\"$bp\" \
    }" > /dev/null
    index=$(($index + 1))
done
