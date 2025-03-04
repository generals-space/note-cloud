# ctr tag重命名镜像名称

```
crictl pull centos:7.2
ctr -n k8s.io image tag centos:7.2 centos:7
```

`crictl images`与`ctr -n k8s.io image ls`的结果一致.
