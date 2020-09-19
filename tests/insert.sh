#!/bin/bash

set -e

if [ $# -ne 1 ]; then
    echo "$0 <count>"
    exit 1
fi

echo "`date`: generating sql..."
index=0
line=""
lbbase="`printf "%x" $(date +%s)`-701c-4170-9073-"
bpbase="`printf "%x" $(date +%s)`-a3b8-4aa4-b559-"
while [ $index -lt $1 ]; do
    # lb=`uuidgen | tr '[:upper:]' '[:lower:]'`
    # bp=`uuidgen | tr '[:upper:]' '[:lower:]'`
    lb=$lbbase`printf "%012x" $index`
    bp=$bpbase`printf "%012x" $index`
    line=$line"insert into mappings(loadbalancer, bigip) values('$lb', '$bp');" 
    
    index=$(($index + 1))

    ready=$(($index % 500))
    if [ $ready -eq 0 ]; then
        echo "`date`: beginning insert..."
        echo "$line" | sqlite3 mapping.db &
        echo "`date`: done $index"
        line=""
    fi
done

echo "`date`: beginning insert..."
echo "$line" | sqlite3 mapping.db &
echo "`date`: done"
