#!/bin/sh
set -eu

commit=$(git rev-parse --short HEAD)
buildDate=$(date -u +'%Y-%m-%dT%H:%M:%SZ')
pkg="github.com/filleokus/hemlis/internal/version"
ldflags="-X '$pkg.Commit=$commit' -X '$pkg.BuildDate=$buildDate'"

go test ./...
go build -v -ldflags="$ldflags" ./cmd/hemlis-combine
go build -v -ldflags="$ldflags" ./cmd/hemlis-gen
