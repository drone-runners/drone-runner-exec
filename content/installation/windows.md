---
date: 2000-01-01T00:00:00+00:00
title: Install on Windows
title_in_header: Windows
author: bradrydzewski
weight: 1
toc: true
description: |
  Install the runner on Windows Server.
---

This article explains how to install the exec runner on Windows. The exec runner is packaged in binary format and distributed as a Github [release](https://github.com/drone-runners/drone-runner-exec/releases).

# Step 1 - Download

Download and unpack the binary.

```
curl -L https://github.com/drone-runners/drone-runner-exec/releases/download/${VERSION}/drone_runner_exec_windows_amd64.tar.gz | tar zx
sudo install -t /usr/local/bin drone-runner-exec
```

# Step 2 - Configure

The exec runner is configured using an environment variable file. The environment file should be stored at the following path: 

```
C:\Drone\drone-runner-exec\config
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

The exec runner writes logs to a file on the host machine. The log file location should be configured in the environment file before you start the service.

```
DRONE_LOG_FILE=C:\Drone\drone-runner-exec\log.txt
```

The log file directory must be created before you start the service:

```
mkdir C:\Drone
mkdir C:\Drone\drone-runner-exec
```

# Step 4 - Install

Install and start the service.

```
drone-runner-exec service create
drone-runner-exec server start
```

# Step 5 - Verify

Inspect the logs and verify the runner successfully established a connection with the Drone server.

```
$ cat C:\Drone\drone-runner-exec\log.txt

INFO[0000] starting the server
INFO[0000] successfully pinged the remote server
```