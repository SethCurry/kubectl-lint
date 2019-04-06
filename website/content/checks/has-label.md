---
title: "has-label"
date: 2019-04-06T17:25:21Z
---

## Explanation

Labels are a very useful utility in Kubernetes for organization, and they are even required for
some types of manifests such as Deployments.

You should almost certainly have at least one label.

## Resolution

Add a label to the pod.  You could use labels to tag pods with:

* The team that runs them
* The version of the software inside them
* The application they are part of, for microservices

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: frontend
  labels:
    owner: Seth
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

