# docker启动失败-kernel overlayfs upper fs needs to support d_type

参考文章

1. [centos7上安装docker 17.06ce，使用的xfs文件系统不支持d-type的问题](https://blog.csdn.net/judyjie/article/details/86178959)
2. [基于XFS文件系统的overlayfs下使用docker，为何要使用d_type=1](https://blog.csdn.net/u014155354/article/details/86648169)

centos: 7
docker: 20.10.12

## 问题描述

在使用kubeadm部署kube集群前, 先安装docker, 但是按照之前文档中修改的`daemon.json`内容, docker却启动失败

```console
$ journalctl -xe
Feb  7 19:28:17 k8s-master-01 kernel: overlayfs: upper fs needs to support d_type.
Feb  7 19:28:17 k8s-master-01 dockerd: failed to start daemon: error initializing graphdriver: overlay2: the backing xfs filesystem is formatted without d_type support, which leads to incorrect behavior. Reformat the filesystem with ftype=1 to enable d_type support. Backing filesystems without d_type support are not supported.
Feb  7 19:28:17 k8s-master-01 systemd: docker.service: main process exited, code=exited, status=1/FAILURE
```

```json
{
    "exec-opts": ["native.cgroupdriver=systemd"],
    "log-driver": "json-file",
    "log-opts": {
        "max-size": "100m"
    },
    "storage-driver": "overlay2",
    "storage-opts": [
        "overlay2.override_kernel_check=true"
    ]
}
```

重启系统也不行.

## 解决方法

按照参考文章1, 2的说法, 内核版本低于4.x的系统, 无法使用overlay2, 即使使用`modprobe overlay`也不行.

于是把"storage-driver"和"storage-opts"去掉了, 再次重启就行了.
