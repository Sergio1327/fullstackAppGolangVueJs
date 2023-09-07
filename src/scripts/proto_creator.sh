#!/bin/bash

if [ -z "$1" ]
  then
    echo "Укажите название модуля"
    exit 1
fi


MODULE=$1
PROTO_PATH=../internal/proto/${MODULE}/
PROTO_FILE_NAME=${MODULE}.proto
PROTO_FILE=${PROTO_PATH}/${PROTO_FILE_NAME}
PROTO_GEN_DIR=${MODULE}proto

if test -f "$PROTO_FILE"; then
    echo "$PROTO_FILE есть, запускаю proto. "
    protoc -I $PROTO_PATH $PROTO_FILE --go_out=plugins=grpc:$PROTO_PATH/
else
    echo "$PROTO_FILE отсутствует, создаю, запустите повторно, чтобы запустить proto генерацию"
    mkdir -p ${PROTO_PATH}
    echo "syntax = \"proto3\";

package proto;
option go_package = './${PROTO_GEN_DIR}';

import \"google/protobuf/timestamp.proto\";

message AuthRequest {
    string login = 1;
    string password = 2;
    bool oauth2 = 3;
}

message AuthResponse {
    Profile profile = 1;
}

message Profile {
    // @inject_tag: db:\"LOGIN\"
    string login = 1;
    // @inject_tag: db:\"SALDO\"
    float balance = 2;
    // @inject_tag: db:\"C_ID\"
    int64 c_id = 3;
    // @inject_tag: db:\"IS_SARKOR_TV\"
    bool is_sarkor_tv = 4;
    bool is_oauth = 5;
    // @inject_tag: db:\"C_DATE\"
    google.protobuf.Timestamp created_date = 6;
}

service TvAPIService {
    rpc Auth(AuthRequest) returns (AuthResponse) {}
}

" > ${PROTO_FILE}
fi

