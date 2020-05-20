#! /bin/sh

rm -rf ./build/*
CGO_ENABLED=0 GOOS=linux ARCH=amd64 go build -ldflags '-extldflags "-static"' -o ./build/nacos-client
cp -r conf ./build
