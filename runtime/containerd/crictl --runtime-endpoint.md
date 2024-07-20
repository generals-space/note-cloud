参考文章

1. [Kubernetes on CRI-O (CentOS)](https://cyylog.netlify.app/2020/12/15/container/kubernetes-on-cri-o-centos/)
    - CN
    - crio服务启动失败及解决方法(LVM): Validating root config: failed to get store to set defaults: kernel does not support overlay fs: overlay: the backing xfs filesystem is formatted without d_type support, which leads to incorrect behavior. Reformat the filesystem with ftype=1 to enable d_type support. Running without d_type is not supported.: driver not supported
2. [Kubernetes on CRI-O (CentOS)](https://dev.to/abhivaidya07/kubernetes-on-cri-o-centos-o1m)
    - EN

```bash
crictl --runtime-endpoint unix:///var/run/containerd/containerd.sock ps -a
```

设置默认值

```bash
echo 'export CONTAINER_RUNTIME_ENDPOINT=unix:///var/run/containerd/containerd.sock' >> ~/.bashrc
source ~/.bashrc
```

永久生效

```
crictl config runtime-endpoint unix:///run/containerd/containerd.sock
```
