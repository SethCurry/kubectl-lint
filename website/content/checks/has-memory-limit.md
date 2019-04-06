---
title: "has-memory-limit"
date: 2019-04-06T17:25:21Z
---

## Explanation

`has-memory-limit` is a warning that a container in one of the pods is missing a memory limit [Resource QoS constraint](https://github.com/kubernetes/community/blob/master/contributors/design-proposals/node/resource-qos.md).
Containers that do not have a memory limit could potentially consume all of the free memory in their node if
they have a memory leak or other form of defect.

## Resolution

The solution is simply to add a memory limit to the manifest.  You can see memory limits in the
manifest below under the `resources` section of the container definitions:

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: frontend
spec:
  containers:
  - name: db
    image: mysql
    env:
    - name: MYSQL_ROOT_PASSWORD
      value: "password"
    resources:
      requests:
        memory: "64Mi"
        cpu: "250m"
      limits:
        memory: "128Mi"
        cpu: "500m"
  - name: wp
    image: wordpress
    resources:
      requests:
        memory: "64Mi"
        cpu: "250m"
      limits:
        memory: "128Mi"
        cpu: "500m"
```

