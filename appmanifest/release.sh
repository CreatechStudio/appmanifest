#!/bin/bash

VERSION="1.0.1"
NAME=appmanifest
OUTPUT=../build

echo "Building $NAME version $VERSION"

mkdir -p ${OUTPUT}

build() {
  echo -n "=> $1-$2: "
  local artifact="${OUTPUT}/${NAME}-$1-$2"
  GOOS=$1 GOARCH=$2 go build -o ${artifact} -ldflags "-X main.version=$VERSION -X main.gitHash=`git rev-parse HEAD`" ./${NAME}.go
  du -h ${artifact}
}

build "darwin" "amd64"
build "darwin" "arm64"
