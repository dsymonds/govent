#! /bin/bash -e

# Requires
#   go install github.com/dsymonds/goembed

cd $(dirname $0)
goembed -package govent -var frontHTML < front.html > front.html.go
gofmt -w -l front.html.go
