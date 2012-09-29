#!/bin/sh

echo "Please make sure, you are currently in the same directory as this script!"
export GOPATH="$GOPATH:$PWD"
echo "GOPATH: $GOPATH"
go install sokoban
echo "done"
