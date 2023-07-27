#!/bin/bash
export ROOT=../../..
source variables.sh
export PG_URL=project_template:project_template@127.0.0.1:5432/project_template
export ORACLE_CONNECT_FILE=../../../../../config/oracle_connect
export DEBUG=true
export DEV=true

go clean -testcache
cd ../internal/
echo '-- тест usecase --'
export CONF_PATH=../../../../config/conf.yaml
for s in $(go list ./usecase/test/...); do if ! go test -failfast -p 1 $s; then break; fi; done 2>&1 | grep -v '/usr/bin/ld: /usr/lib/oracle/11.2/client64/lib/libnnz11.so'
echo '-- тест repository --'
export CONF_PATH=../../../../../config/conf.yaml
for s in $(go list ./repository/oracle/test/...); do if ! go test -failfast -p 1 $s; then break; fi; done 2>&1 | grep -v '/usr/bin/ld: /usr/lib/oracle/11.2/client64/lib/libnnz11.so'
for s in $(go list ./repository/postgresql/test/...); do if ! go test -failfast -p 1 $s; then break; fi; done 2>&1 | grep -v '/usr/bin/ld: /usr/lib/oracle/11.2/client64/lib/libnnz11.so'