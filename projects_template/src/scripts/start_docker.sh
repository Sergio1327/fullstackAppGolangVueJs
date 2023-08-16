#!/usr/bin/env bash

export ROOT=../../docker
source variables.sh

echo 'RUN DOCKER'

cd $ROOT
docker-compose build
docker-compose up --force-recreate --remove-orphans


