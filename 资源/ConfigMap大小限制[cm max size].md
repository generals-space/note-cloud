# kuber-ConfigMap大小限制[cm max size]

参考文章

1. [Kubernetes ConfigMap size limitation](https://stackoverflow.com/questions/53012798/kubernetes-configmap-size-limitation)
2. [Size limit for ConfigMap](https://github.com/kubernetes/kubernetes/issues/19781)
    - 参考文章1引用了 issue

`ConfigMap`本身没有大小限制, 但是`ConfigMap`需要存储在`etcd`中, `etcd`对每个key的value大小是有限制的, 最大为1M.

