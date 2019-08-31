---
date: 2000-01-01T00:00:00+00:00
title: Install on OSX
title_in_header: OSX
author: bradrydzewski
weight: 1
toc: true
description: |
  Install the runner on OSX.
---

This article explains how to install the exec runner on OSX. The exec runner is packaged in binary format and distributed as a Github [release](https://github.com/drone-runners/drone-runner-exec/releases).

# Step 1 - Download

Download and unpack the binary. You will need to modify the below commands to use the latest release.

```text
$ curl -L https://github.com/drone-runners/drone-runner-exec/releases/download/${VERSION}/drone_runner_exec_darwin_amd64.tar.gz | tar zx
sudo cp drone-runner-exec /usr/local/bin
```

# Step 2 - Configure

The exec runner is configured using an environment variable file. The environment file location is determined by the user running the process. If the user is root: 

```
/etc/drone-runner-exec/config
```

If the user is not-root:

```
~/.drone-runner-exec/config
```

This file should use the syntax `<variable>=value` which sets the variable to the given value and `#` for comments. Please note this is not a bash file. Bash syntax and Bash expressions are not supported.

```
DRONE_RPC_PROTO=https
DRONE_RPC_HOST=drone.company.com
DRONE_RPC_SECRET=super-duper-secret
```

This article references the below configuration options. See [Configuration]({{< relref "reference" >}}) for a complete list of configuration options.

DRONE_RPC_HOST
: provides the hostname (and optional port) of your Drone server. The runner connects to the server at the host address to receive pipelines for execution.

DRONE_RPC_PROTO
: provides the protocol used to connect to your Drone server. The value must be either http or https.

DRONE_RPC_SECRET
: provides the shared secret used to authenticate with your Drone server. This must match the secret defined in your Drone server configuration.

# Step 3 - Logging

The exec runner writes logs to a file on the host machine. The log file location should be configured in the environment file before you start the service. The file location is determined by the user running the process. If the user is root: 

```
DRONE_LOG_FILE=/var/log/drone-runner-exec/log.txt
```

If the user is non-root:

```
DRONE_LOG_FILE=.drone-runner-exec/log.txt
```

The log file directory must be created before you start the service:

```
$ mkdir /var/log/drone-runner-exec
$ mkdir .drone-runner-exec/
```

# Step 4 - Install

Install and start the service.

```
$ drone-runner-exec service install
$ drone-runner-exec server start
```

# Step 5 - Verify

Inspect the logs and verify the runner successfully established a connection with the Drone server. If the user is root:

```
$ cat /var/log/drone-runner-exec/log.txt

INFO[0000] starting the server
INFO[0000] successfully pinged the remote server
```

If the user is non-root:

```
$ cat ~/.drone-runner-exec/log.txt

INFO[0000] starting the server
INFO[0000] successfully pinged the remote server
```
