#!/bin/sh
# go test -v ./... 2>&1 | sed "s/.* assertion.*/[0m/" | grep -v -E "^(-|$|\?|PASS|ok)"

# go test -v ./... -run CreateListMsg
go test ./...  | grep -iEv "^\?"
