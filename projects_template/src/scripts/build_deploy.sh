#!/usr/bin/env bash

source variables.sh

export ROOT=../../docker
export IMAGE="$NAME:$VERSION"

echo 'build docker'
cd $ROOT
docker compose build 

docker tag $IMAGE dock.sk.uz/$IMAGE
docker push dock.sk.uz/$IMAGE
