---
date: 2000-01-01T00:00:00+00:00
title: Platform
author: bradrydzewski
weight: 2
toc: false
description: |
  Configure the target operating system and architecture.
---

Use the `platform` section to configure the target operating system and architecture and route the pipeline to the appropriate runner.

Example macOS (darwin) pipeline:

{{< highlight text "linenos=table,hl_lines=5-7" >}}
kind: pipeline
type: exec
name: default

platform:
  os: darwin
  arch: amd64

steps:
- name: build
  commands:
  - go build
  - go test
{{< / highlight >}}

# Supported Platforms

os          | arch
------------|-----
`linux`     | `amd64`
`linux`     | `arm64`
`linux`     | `arm`
`linux`     | `386`
`windows`   | `amd64`
`windows`   | `386`
`darwin`    | `amd64`
`freebsd`   | `amd64`
`freebsd`   | `arm`
`freebsd`   | `386`
`netbsd`    | `amd64`
`netbsd`    | `arm`
`openbsd`   | `amd64`
`openbsd`   | `arm`
`openbsd`   | `386`
`dragonfly` | `amd64`
`solaris`   | `amd64`



<!-- # Linux

Linux operating system, amd64 architecture:

{{< highlight text "linenos=table,linenostart=5" >}}
platform:
  os: linux
  arch: amd64
{{< / highlight >}}

Linux operating system, i386 architecture:

{{< highlight text "linenos=table,linenostart=5" >}}
platform:
  os: linux
  arch: 386
{{< / highlight >}}

Linux operating system, arm64 architecture:

{{< highlight text "linenos=table,linenostart=5" >}}
platform:
  os: linux
  arch: arm64
{{< / highlight >}}

Linux operating system, arm architecture:

{{< highlight text "linenos=table,linenostart=5" >}}
platform:
  os: linux
  arch: arm
{{< / highlight >}}

# Windows

Windows operating system, amd64 architecture:

{{< highlight text "linenos=table,linenostart=5" >}}
platform:
  os: windows
  arch: amd64
{{< / highlight >}}

Windows operating system, i386 architecture:

{{< highlight text "linenos=table,linenostart=5" >}}
platform:
  os: windows
  arch: 386
{{< / highlight >}}

# Darwin

Darwin (OSX) operating system, amd64 architecture:

{{< highlight text "linenos=table,linenostart=5" >}}
platform:
  os: darwin
  arch: amd64
{{< / highlight >}}

# FreeBSD

FreeBSD operating system, amd64 architecture:

{{< highlight text "linenos=table,linenostart=5" >}}
platform:
  os: freebsd
  arch: amd64
{{< / highlight >}}

FreeBSD operating system, i386 architecture:

{{< highlight text "linenos=table,linenostart=5" >}}
platform:
  os: freebsd
  arch: 386
{{< / highlight >}}

FreeBSD operating system, arm architecture:

{{< highlight text "linenos=table,linenostart=5" >}}
platform:
  os: freebsd
  arch: arm
{{< / highlight >}} -->

