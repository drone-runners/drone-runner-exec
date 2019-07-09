#!/bin/sh

set -e
set -x

# linux
tar -cvzf release/drone_runner_exec_linux_amd64.tar.gz -C release/linux/amd64 drone-runner-exec
tar -cvzf release/drone_runner_exec_linux_arm64.tar.gz -C release/linux/arm64 drone-runner-exec
tar -cvzf release/drone_runner_exec_linux_arm.tar.gz   -C release/linux/arm   drone-runner-exec
tar -cvzf release/drone_runner_exec_linux_386.tar.gz   -C release/linux/386   drone-runner-exec

# windows
tar -cvzf release/drone_runner_exec_windows_amd64.tar.gz -C release/windows/amd64 drone-runner-exec
tar -cvzf release/drone_runner_exec_windows_386.tar.gz   -C release/windows/386   drone-runner-exec

# darwin
tar -cvzf release/drone_runner_exec_darwin_amd64.tar.gz -C release/darwin/amd64  drone-runner-exec

# freebase
tar -cvzf release/drone_runner_exec_freebsd_amd64.tar.gz -C release/freebsd/amd64 drone-runner-exec
tar -cvzf release/drone_runner_exec_freebsd_arm.tar.gz   -C release/freebsd/arm   drone-runner-exec
tar -cvzf release/drone_runner_exec_freebsd_386.tar.gz   -C release/freebsd/386   drone-runner-exec

# netbsd
tar -cvzf release/drone_runner_exec_netbsd_amd64.tar.gz -C release/netbsd/amd64 drone-runner-exec
tar -cvzf release/drone_runner_exec_netbsd_arm.tar.gz   -C release/netbsd/arm   drone-runner-exec

# openbsd
tar -cvzf release/drone_runner_exec_openbsd_amd64.tar.gz -C release/openbsd/amd64 drone-runner-exec
tar -cvzf release/drone_runner_exec_openbsd_arm.tar.gz   -C release/openbsd/arm   drone-runner-exec
tar -cvzf release/drone_runner_exec_openbsd_386.tar.gz   -C release/openbsd/386   drone-runner-exec

# dragonfly
tar -cvzf release/drone_runner_exec_dragonfly_amd64.tar.gz -C release/dragonfly/amd64  drone-runner-exec

# solaris
tar -cvzf release/drone_runner_exec_solaris_amd64.tar.gz -C release/solaris/amd64  drone-runner-exec

# generate shas for tar files
shasum release/*.tar.gz > release/drone_runner_exec_checksums.txt
