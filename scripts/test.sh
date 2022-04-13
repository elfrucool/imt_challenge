#!/bin/sh

go list ./... | while read LINE ; do echo "===[$LINE]==="; go test $LINE "$@" ; done
