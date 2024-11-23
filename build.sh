#!/bin/bash
rm -rf ./build/*
go build -o ./build/task-cli task-cli.go
