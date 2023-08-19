#!/bin/bash

#build project
rm feishusign
rm -rf ./logs
go build -o feishusign ./cmd/main.go
