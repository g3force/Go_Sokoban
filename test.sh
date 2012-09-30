#!/bin/sh

echo "Please make sure, you are currently in the same directory as this script!"
export GOPATH="$GOPATH:$PWD"
echo "GOPATH: $GOPATH"
go test sokoban/ai
go test sokoban/engine
echo "done"
