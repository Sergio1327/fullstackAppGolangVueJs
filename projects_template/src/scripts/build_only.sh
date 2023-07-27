#!/usr/bin/env bash

source variables.sh

export BUILD_TIME=`date  +%d.%m.%Y-%H:%M:%S`
export LAST_COMMIT=$(git rev-parse --verify HEAD)
LAST_COMMIT=${LAST_COMMIT:0:8}

mkdir -p $ROOT/log
mkdir -p $ROOT/bin

cd $ROOT/cmd
echo 'COMPILING...'
rm $ROOT/bin/$BIN
go build -o $BIN -ldflags "-X main.version=$BUILD_TIME|$LAST_COMMIT" || { echo 'build failed' ; exit 1; }
mv $BIN $ROOT/bin
