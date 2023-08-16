#!/bin/bash

export DB_NAME=test_db
export DB_PASSWORD=test_db
ROOT=../..


migrate -path $ROOT/docker/migrate -database postgres://$DB_NAME:$DB_PASSWORD@localhost:5432/$DB_NAME?sslmode=disable up 1