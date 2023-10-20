#!/bin/bash

mkdir -p bin/
rm -rf bin/*
cd bin/

go build ../cmd/sitesync/
go build ../cmd/deploy/
go build ../cmd/kubernetes

cd ..
