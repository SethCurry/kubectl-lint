---
title: "no-default-namespace"
date: 2019-04-06T17:25:21Z
---

## Explanation

Namespaces are a very useful tool for logically separating resources.

It's generally discouraged to deploy everything into the `default` namespace.

## Resolution

Move the object to a different namespace.  If you do not have another namespace, you can create one with:

```bash
kubectl create namespace <namespace-name>
```

And then change your namespace in the object manifest using the `metadata.namespace` field:

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: frontend
  namespace: <namespace-name>
spec:
  containers:
  - name: db
    image: mysql
```

