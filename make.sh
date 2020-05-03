#!/bin/bash
mkdir -p /tmp/gomics
sudo cp -r "$PWD/." /tmp/gomics
odir="$PWD"
cd /tmp/gomics/src
GOPATH="$PWD" # for dependencies
go generate
go build
mkdir -p ../bin
mv src ../bin/gomics
cd "$odir"
cp /tmp/gomics/bin/gomics "$PWD/bin"
