---
title: "Configuration File"
date: 2019-04-06T02:58:42Z
---

You can configure `kubectl-lint` by using a configuration file.
By default, it will check for configuration files at:

* ./.kubectl-lint.json
* ./.kubectl-lint.yaml
* $HOME/.kubectl-lint.json
* $HOME/.kubectl-lint.yaml

## Example Configuration File

```yaml
disable:
  - no-cpu-limit
  - no-memory-limit
```

