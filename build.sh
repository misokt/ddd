#!/bin/bash

# 64-bit linux
GOOS=linux GOARCH=amd64 go build
zip linux-64bit-dump-discord-data.zip dump-discord-data
echo "finished 64-bit linux"
rm dump-discord-data

# 32-bit linux
GOOS=linux GOARCH=386 go build
zip linux-32bit-dump-discord-data.zip dump-discord-data
echo "finished 32-bit linux"
rm dump-discord-data

# 64-bit macos
GOOS=darwin GOARCH=amd64 go build
zip macos-64bit-dump-discord-data.zip dump-discord-data
echo "finished 64-bit macos"
rm dump-discord-data

# silicon macos
GOOS=darwin GOARCH=arm64 go build
zip macos-silicon-dump-discord-data.zip dump-discord-data
echo "finished silicon macos"
rm dump-discord-data

# 64-bit windows
GOOS=windows GOARCH=amd64 go build
zip windows-64bit-dump-discord-data.zip dump-discord-data.exe
echo "finished 64-bit windows"
rm dump-discord-data.exe

# 64-bit windows
GOOS=windows GOARCH=386 go build
zip windows-32bit-dump-discord-data.zip dump-discord-data.exe
echo "finished 64-bit windows"
rm dump-discord-data.exe
