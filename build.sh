#!/bin/zsh

GOOS=darwin GOARCH=amd64 go build -o meton_darwin-amd64 main.go
GOOS=linux GOARCH=amd64 go build -o meton_linux-amd64 main.go
# GOOS=windows GOARCH=amd64 go build -o meton_windows-amd64 main.go