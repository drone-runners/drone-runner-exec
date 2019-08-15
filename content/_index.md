---
date: 2000-01-01T00:00:00+00:00
title: Exec Runner
author: bradrydzewski
weight: 1
expand: configuration/_index.md
---

The exec runner is a daemon (aka agent) that executes build pipelines directly on the host machine without isolation. This documentation provides details for installing, configuring and using the exec runner.

{{< alert "security" >}}
This runner executes commands directly on the host machine. It does not provide isolation or defense against malicious code. This runner should only be used in a trusted environment.
{{< / alert >}}

If you want to install this runner:

{{< link "/installation" >}}

If you want to configure this runner for your project:

{{< link "/configuration" >}}

If you have questions or require assistance:

{{< link "/support" >}}