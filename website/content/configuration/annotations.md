---
title: "Annotations"
date: 2019-04-06T02:58:42Z
---

You can also use Kubernetes Annotations to alter the behavior of `kubectl-lint`
for a single object/manifest.

## Current Annotations

```yaml
metadata:
  annotations:
    lint/disable-linters: "linter1,linter2"
```

