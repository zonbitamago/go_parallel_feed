#!/bin/sh

go test -v -coverprofile cover.out ./...
go tool cover -html=cover.out -o cover.html