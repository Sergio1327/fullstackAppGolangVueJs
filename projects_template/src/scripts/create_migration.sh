#!/bin/bash

if [ -z "$1" ]
  then
    echo "Укажите название миграции"
    exit 1
fi

NAME=$1
ROOT=../..

migrate create -ext sql -dir $ROOT/docker/migrate -tz=Asia/Tashkent -seq  $NAME