---
title: "has-memory-request"
date: 2019-04-06T17:25:21Z
---

## Explanation

`has-memory-request` is a warning that a container in one of the pods is missing a memory request [Resource QoS constraint](https://github.com/kubernetes/community/blob/master/contributors/design-proposals/node/resource-qos.md).

The pod scheduler will not be able to determine how much memory your container will use, which can result in more Out of Memory issues.
Adding a memory request will give the scheduler a hint about the minimum amount of memory your container will need, which should lead
to better scheduling.

## Resolution

The solution is simply to add a memory request to the manifest.  You can see memory requests in the
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

