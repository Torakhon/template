#!/bin/bash
CURRENT_DIR=$(pwd)
export GOPATH=~/projects/go
export PATH=$PATH:/$GOPATH/bin

# shellcheck disable=SC2044
for module in $(find "$CURRENT_DIR"/protos/* -type d); do
    protoc -I /usr/local/include \
           -I $GOPATH/pkg/mod/github.com/gogo/protobuf \
           -I "$CURRENT_DIR"/protos/ \
            --gofast_out=plugins=grpc:"$CURRENT_DIR"/genproto/ \
            "$module"/*.proto;
done;

# shellcheck disable=SC2044
for module in $(find "$CURRENT_DIR"/genproto/* -type d); do
  if [[ "$OSTYPE" == "darwin"* ]]; then
    sed -i "" -e "s/,omitempty//g" "$module"/*.go
  else
    # shellcheck disable=SC2086
    sed -i -e "s/,omitempty//g" $module/*.go
  fi
done;