#! /bin/bash -e

# Requires
#   go install github.com/dsymonds/goembed

goembed -package govent -var frontHTML < front.html > front.html.go
gofmt -w -l front.html.go
