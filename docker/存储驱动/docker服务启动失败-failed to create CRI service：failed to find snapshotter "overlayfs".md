参考文章

1. [container的构建镜像失败：snapshotter not loaded: overlayfs: invalid argument](https://blog.csdn.net/jieshibendan/article/details/122574854)
2. [Use the OverlayFS storage driver](https://docs.docker.com/storage/storagedriver/overlayfs-driver/#prerequisites)
    - `overlay`和`overlay2`可以在`xfs`文件系统上运行, 但要求其拥有`d_type=true`配置
    - `xfs_info | grep ftype`, 可以查看当前文件系统是否支持`d_type=true`, 1表示支持`d_type`, 0表示不支持
    - 不过, 要想让一个磁盘分区支持`d_type`, 只能在格式化分区时操作, 如`mkfs.xfs -n ftype=1`, 没有办法对一个正在使用的分区进行转换.
3. [基于XFS文件系统的OverlayFS](https://blog.csdn.net/avatar_2009/article/details/107666571)
    - xfs文件系统的`d_type`含义

centos 7: 3.10.0-1062.4.1.el7.x86_64


`systemctl start docker`无法启动docker服务, 查看`/var/log/message`中, 发现有问题的只有如下日志.

```log
May 15 10:13:50 k8s-master-01 containerd: time="2023-05-15T10:13:50.751476886+08:00" level=info msg="Connect containerd service"
May 15 10:13:50 k8s-master-01 containerd: time="2023-05-15T10:13:50.751857807+08:00" level=warning msg="failed to load plugin io.containerd.grpc.v1.cri" error="failed to create CRI service: failed to find snapshotter \"overlayfs\""
May 15 10:13:50 k8s-master-01 containerd: time="2023-05-15T10:13:50.751889826+08:00" level=info msg="loading plugin \"io.containerd.grpc.v1.introspection\"..." type=io.containerd.grpc.v1
May 15 10:13:50 k8s-master-01 containerd: time="2023-05-15T10:13:50.752220794+08:00" level=info msg=serving... address=/run/containerd/containerd.sock.ttrpc
May 15 10:13:50 k8s-master-01 containerd: time="2023-05-15T10:13:50.752266446+08:00" level=info msg=serving... address=/run/containerd/containerd.sock
May 15 10:13:50 k8s-master-01 containerd: time="2023-05-15T10:13:50.753873854+08:00" level=info msg="containerd successfully booted in 0.090845s"
```

不过`containerd`服务是正常的...

```console
$ ps -ef | grep containerd
root       2161      1  0 20:47 ?        00:00:00 /usr/local/bin/containerd
```

按照参考文章1中所说

```console
$ ctr plugins list
TYPE                            ID                       PLATFORMS      STATUS
io.containerd.snapshotter.v1    aufs                     linux/amd64    skip
io.containerd.snapshotter.v1    btrfs                    linux/amd64    skip
io.containerd.snapshotter.v1    devmapper                linux/amd64    error
io.containerd.snapshotter.v1    native                   linux/amd64    ok
io.containerd.snapshotter.v1    overlayfs                linux/amd64    error
io.containerd.snapshotter.v1    zfs                      linux/amd64    skip
```

`overlayfs`插件加载失败, 而 docker 默认需要`overlayfs`驱动, 因此启动失败.

```json
// /etc/docker/daemon.json
{
  "storage-driver": "overlay2"
}
```


但使用`lsmod`可以看到该内核模块已经加载完成了.

```
$ lsmod | grep overlay
overlay                91659  0
```
