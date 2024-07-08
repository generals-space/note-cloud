# kubeadm init失败 - apiserver 日志报错x509：certificate signed by unknown authority

参考文章

1. [Kubelet client certificate rotation fails](https://kubernetes.io/docs/setup/production-environment/tools/kubeadm/troubleshooting-kubeadm/)

## 问题描述

ubuntu: 20.04

kubernetes: v1.25.1

使用`kubeadm init`初始化 kube 集群失败, kubelet已经启动, `crictl ps`也可以看到各组件已经处于Running状态.

查看系统日志, 发现 kubelet 还存在如下输出.

```log
$ tail -f /var/log/syslog
Sep 22 19:15:01 ubuntu kubelet[10156]: E0922 19:15:01.071785   10156 kubelet.go:2448] "Error getting node" err="node "ubuntu" not found"
Sep 22 19:15:01 ubuntu kubelet[10156]: E0922 19:15:01.172783   10156 kubelet.go:2448] "Error getting node" err="node "ubuntu" not found"
Sep 22 19:15:01 ubuntu kubelet[10156]: E0922 19:15:01.275082   10156 kubelet.go:2448] "Error getting node" err="node "ubuntu" not found"
```

由于之前的问题, 我们知道这是因为 kubelet 请求 apiserver 失败(不过这次已经不再是"connection refused"), 于是查看 apiserver 容器的日志, 看看有什么问题.

```log
$ crictl logs -f apiserver-xxx
E0923 02:12:23.620666       1 authentication.go:63] "Unable to authenticate the request" err="[x509: certificate signed by unknown authority, verifying certificate SN=1707380803054755332, SKID=, AKID=0B:15:61:12:EE:8E:8B:DF:E3:41:59:5A:67:2A:5C:11:02:AE:71:4D failed: x509: certificate signed by unknown authority (possibly because of \"crypto/rsa: verification error\" while trying to verify candidate authority certificate \"kubernetes\")]"
```

看来是 kubelet 认证的问题.

之前由于 kubeadm init 未完成, 所以 kubeadm reset 也总是失败, 应该是 kubelet 在某处的缓存未清理.

按照参考文章1中所说, 手动清理数据时, 要删除如下3个目录

- /var/lib/kubelet
- /var/lib/etcd
- /etc/kubernetes

重新init, 就可以了.
