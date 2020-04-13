# ephemeral container临时容器

参考文章

1. [临时容器](https://kubernetes.io/zh/docs/concepts/workloads/pods/ephemeral-containers/)

kuber版本: 1.16.2, 还是 Alpha.

```
$ k exec coredns-67c766df46-cfdh2 /bin/sh
OCI runtime exec failed: exec failed: container_linux.go:346: starting container process caused "exec: \"/bin/sh\": stat /bin/sh: no such file or directory": unknown
command terminated with exit code 126
```

