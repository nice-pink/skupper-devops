#!/bin/bash

mkdir -p bin/ && cd bin/

go build ../cmd/sitesync/ ../cmd/deploy/ ../cmd/kubernetes

cd ..
