#!/bin/sh

# disable go modules
export GOPATH=""

# disable cgo
export CGO_ENABLED=0

set -e
set -x

# linux
GOOS=linux GOARCH=amd64 go build -o release/linux/amd64/drone-runner-exec
GOOS=linux GOARCH=arm64 go build -o release/linux/arm64/drone-runner-exec
GOOS=linux GOARCH=arm   go build -o release/linux/arm/drone-runner-exec
GOOS=linux GOARCH=386   go build -o release/linux/386/drone-runner-exec

# windows
GOOS=windows GOARCH=amd64 go build -o release/windows/amd64/drone-runner-exec.exe
GOOS=windows GOARCH=386   go build -o release/windows/386/drone-runner-exec.exe

# darwin
GOOS=darwin GOARCH=amd64 go build -o release/darwin/amd64/drone-runner-exec
GOOS=darwin GOARCH=arm64 go build -o release/darwin/arm64/drone-runner-exec

# freebsd
GOOS=freebsd GOARCH=amd64 go build -o release/freebsd/amd64/drone-runner-exec
GOOS=freebsd GOARCH=arm   go build -o release/freebsd/arm/drone-runner-exec
GOOS=freebsd GOARCH=386   go build -o release/freebsd/386/drone-runner-exec

# netbsd
GOOS=netbsd GOARCH=amd64 go build -o release/netbsd/amd64/drone-runner-exec
GOOS=netbsd GOARCH=arm   go build -o release/netbsd/arm/drone-runner-exec

# openbsd
GOOS=openbsd GOARCH=amd64 go build -o release/openbsd/amd64/drone-runner-exec
GOOS=openbsd GOARCH=arm   go build -o release/openbsd/arm/drone-runner-exec
GOOS=openbsd GOARCH=386   go build -o release/openbsd/386/drone-runner-exec

# dragonfly
GOOS=dragonfly GOARCH=amd64 go build -o release/dragonfly/amd64/drone-runner-exec

# solaris
GOOS=solaris GOARCH=amd64 go build -o release/solaris/amd64/drone-runner-exec
