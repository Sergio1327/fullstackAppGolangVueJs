#!/bin/bash
export DEBUG=true
export CONF_PATH=../../config/conf.yaml
export LOG_DIR=../../log
export MEASURE=enable
export ORACLE_CONNECT_FILE=../../config/oracle_connect



if [ -z "$1" ]
  then
    echo "Укажите название tools"
    exit 1
fi

MODULE=$1
FUNC_NAME=""

if [ -n "$2" ]
  then
    FUNC_NAME="-v --run $2"
    echo $FUNC_NAME
fi

cd ../tools/$MODULE
go test -bench=. $FUNC_NAME