#!/usr/bin/env bash

mkdir -p bin
GOBIN=`pwd`/bin go install ./...