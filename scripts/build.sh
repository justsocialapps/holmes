#!/bin/sh

VERSION=$1

GOOS=windows GOARCH=amd64 go build -o holmes_${VERSION}_windows_amd64
GOOS=windows GOARCH=386 go build -o holmes_${VERSION}_windows_386
GOOS=linux GOARCH=386 go build -o holmes_${VERSION}_linux_386
GOOS=linux GOARCH=amd64 go build -o holmes_${VERSION}_linux_amd64
