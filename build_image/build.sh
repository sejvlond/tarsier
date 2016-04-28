#!/usr/bin/env bash

set -e

if [ -z "$1" ]; then
    echo "Provide project name"
    exit 1
fi

PROJECT=$GOPATH/src/$1

mkdir -p $PROJECT
cd $_
cp -R /src/* .
rm -rf vendor
glide install
go test `glide novendor`
go build -tags=netgo -o /src/`basename $PROJECT`
