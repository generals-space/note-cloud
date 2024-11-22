# kubeadm部署后没有容器启动[containerd]

参考文章

1. [Impossible to create or start a container after reboot (OCI runtime create failed: expected cgroupsPath to be of format \"slice:prefix:name\" for systemd cgroups, got \"/kubepods/burstable/...") #4857](https://github.com/containerd/containerd/issues/4857)

- centos: 7
- docker: 19.03.5
- containerd: 1.6.33
- kube: 1.17.2

## 问题描述

kubeadm init 后没有容器启动, crictl ps -a 没有任何容器.

查看`/var/log/messages`, 一直在刷如下日志

```log
Nov 22 16:53:37 k8s-master-01 kubelet: E1122 16:53:37.438501   22530 reflector.go:153] k8s.io/kubernetes/pkg/kubelet/kubelet.go:458: Failed to list *v1.Node: Get https://kube-apiserver.generals.space:6443/api/v1/nodes?fieldSelector=metadata.name%3Dk8s-master-01&limit=500&resourceVersion=0: dial tcp 10.0.2.20:6443: connect: connection refused
Nov 22 16:53:37 k8s-master-01 kubelet: E1122 16:53:37.513866   22530 kubelet.go:2263] node "k8s-master-01" not found
Nov 22 16:53:37 k8s-master-01 kubelet: E1122 16:53:37.615553   22530 kubelet.go:2263] node "k8s-master-01" not found
```

但是我已经将`k8s-master-01`和`kube-apiserver.generals.space`配到了`/etc/hosts`里, 可以ping通, 所以根本原因不是这个, 而是因为 apiserver 没启动.

由于 kubelet 是通过 containerd 创建容器的, 所以过滤`cat /var/log/messages | grep containerd:`发现有如下日志.

```log
Nov 22 16:42:12 k8s-master-01 containerd: time="2024-11-22T16:42:12.359539679+08:00" level=error msg="RunPodSandbox for &PodSandboxMetadata{Name:kube-apiserver-k8s-master-01,Uid:af685e7d4585dfcacbcd21ee025739b6,Namespace:kube-system,Attempt:0,} failed, error" error="failed to create containerd task: failed to create shim task: OCI runtime create failed: runc create failed: expected cgroupsPath to be of format \"slice:prefix:name\" for systemd cgroups, got \"/kubepods/burstable/podaf685e7d4585dfcacbcd21ee025739b6/2f4a2a37a4d4f7f6911373e9e24b757f2c0e02904c4eeb144ac829a7a9768cfe\" instead: unknown"
```

搜索了一下, 按照参考文章1中所说, 这个问题是因为 containerd 默认的 cgroup 管理器为 systemd, 而 kubelet 默认为 cgroupfs.

## 解决方案

修改`/etc/containerd/config.toml`文件, 搜索`SystemdCgroup = true`, 将其改为 false, 然后重启 containerd 即可.

------

不过我构建 k8s 集群, containerd 和 kubelet 关于 cgroup 的配置一直都是默认的, 为啥突然报这个错...
