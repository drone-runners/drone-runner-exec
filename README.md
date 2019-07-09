# drone-runner-exec

The `exec` runner executes pipelines directly on the host machine. This runner is intended for workloads that are not suitable for running inside containers. This requires Drone server `1.2.1` or higher.

__Warning__ that the exec runner should only be used in a trusted environment. The exec runner grants the pipeline full access to the host machine. There is no isolation.

_Please note the exec runner is experimental and may not be suitable for production use. Furthermore, this runner is not yet subject to support or service level agreements._

# Todos

- [ ] detailed documentation at docs.drone.io
- [x] support for client side filtering
- [x] support for detached steps
- [ ] test windows
- [ ] test launchd file in contrib for osx
- [ ] test systemd file in contrib for linux
- [ ] provide windows service configuration file

# Installation

## Download

Download and install the runner:

```cmd
$ curl -L https://github.com/drone-runners/drone-runner-exec/releases/download/v1.0.0-beta.1/drone_runner_exec_darwin_amd64.tar.gz | tar zx
$ sudo cp drone-runner-exec /usr/local/bin
```

_Please note that you must download the binary distribution that matches your target operating system and architecture. You should adjust the above command accordingly._

## Configure and Run

```cmd
export DRONE_RPC_PROTO=http
export DRONE_RPC_HOST=drone.company.com
export DRONE_RPC_SECRET=super-duper-secret
/usr/local/bin/drone-runner-exec
```

## Debugging

Set the following variable to debug generic execution issues:

```sh
DRONE_DEBUG=true
DRONE_TRACE=true
```

Set the following variables to debug server communication issues:

```sh
DRONE_DEBUG=true
DRONE_TRACE=true
DRONE_RPC_DUMP_HTTP=true
DRONE_RPC_DUMP_HTTP_BODY=true
```

# Usage

## Definition

Use the `kind` and `type` attributes to indicate the pipeline should be routed to the exec runner.

```yaml
kind: pipeline
type: exec
name: default
```

## Platforms

Use the `platform` block to define the target platform and route the build to the appropriate runner. The platform should always be defined.

```yaml
kind: pipeline
type: exec
name: default

platform:
  os: darwin
  arch: amd64
```

List of supported platforms:

os          | arch
------------|-----
`linux`     | `amd64`
`linux`     | `amd64`
`linux`     | `arm`
`linux`     | `arm64`
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

## Steps

Pipeline steps are defined in the `steps` block and are executed serially by default. You can use the `when` clause to limit step execution at runtime.

```yaml
kind: pipeline
type: exec
name: default

platform:
  os: darwin
  arch: amd64

steps:
- name: test
  commands:
  - go test

- name: build
  commands:
  - go build
  when:
    event: [ push ]
```

_Please note the `exec` pipeline configuration is a distinct configuration format. It shares some similarities with the docker pipeline configuration._
