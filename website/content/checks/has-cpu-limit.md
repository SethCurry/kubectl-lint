---
title: "has-cpu-limit"
date: 2019-04-06T17:25:21Z
---

## Explanation

`has-cpu-limit` is a warning that a container in one of the pods is missing a CPU limit [Resource QoS constraint](https://github.com/kubernetes/community/blob/master/contributors/design-proposals/node/resource-qos.md).
Containers that do not have a memory CPU could potentially consume all of the free CPU in their node if
they have a deadlock or other form of CPU-consuming bug.

## Resolution

The solution is simply to add a CPU limit to the manifest.  You can see CPU limits in the
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

