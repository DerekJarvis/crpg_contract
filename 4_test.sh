#!/bin/bash
go clean -testcache && go test ./chaincode -p 1 -v
