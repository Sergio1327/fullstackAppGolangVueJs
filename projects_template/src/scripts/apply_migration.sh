#!/bin/bash

export DB_NAME=project_template
export DB_PASSWORD=project_template

migrate -path ../docker/migrate -database postgres://$DB_NAME:$DB_PASSWORD@localhost:5432/$DB_NAME?sslmode=disable up 1