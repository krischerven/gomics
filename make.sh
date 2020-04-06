#!/bin/bash
cd src
go generate
go build
mkdir -p ../bin
mv src ../bin/gomics
