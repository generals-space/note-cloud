参考文章

1. [Docker已经再见，替代 Docker 的五种容器选择](https://cloud.tencent.com/developer/article/1422822)
    - [apache/mesos](https://github.com/apache/mesos) C++
    - [rkt/rkt](https://github.com/rkt/rkt) 该项目已结束, CoreOS公司最早发起.
    - [docker/docker-ce](https://github.com/docker/docker-ce)
    - LXC 容器. 不支持与 kuber 整合, 没有实现 OCI 的标准

[opencontainers/runc](https://github.com/opencontainers/runc)
    - 之前 docker 旗下的 [libcontainer](https://github.com/docker-archive/libcontainer)

[containerd/containerd](https://github.com/containerd/containerd)
    - 之前 docker 旗下的 [containerd](https://github.com/docker-archive/containerd)

docker info: 19.03.5

相关的可执行文件有:

```console
$ ls /usr/bin/ | grep docker
docker
dockerd
docker-init
docker-proxy
$ ls /usr/bin/ | grep container
containerd
containerd-shim
$ ls /usr/bin/ | grep runc
runc
```

```console
$ ps -ef | grep containerd
root       1258      1  0 05:58 ?        00:00:44 /usr/bin/containerd
root       1263      1  1 05:58 ?        00:06:17 /usr/bin/dockerd -H fd:// --containerd=/run/containerd/containerd.sock
root       1942   1258  0 05:58 ?        00:00:04 containerd-shim -namespace moby -workdir /var/lib/containerd/io.containerd.runtime.v1.linux/moby/ad7a7a5952fd0f1b6637d49cdb673d73b73b65d750f21a734b133d8e07e25b98 -address /run/containerd/containerd.sock -containerd-binary /usr/bin/containerd -runtime-root /var/run/docker/runtime-runc -systemd-cgroup
```

可以看到, `dockerd`与`containerd`是并列的, `containerd-shim`(每启动一个docker容器都会启动一个shim进程)是`containerd`的子进程, 容器中的`CMD/ENTRYPOINT`执行命令是由`containerd-shim`启动执行的. 

如下, 容器`ad7a7a5952fd0`中的nginx进程就是对应`container-shim`的子进程.

```console
$ ps -ef | grep 1942
root       1942   1258  0 05:58 ?        00:00:04 containerd-shim -namespace moby -workdir /var/lib/containerd/io.containerd.runtime.v1.linux/moby/ad7a7a5952fd0f1b6637d49cdb673d73b73b65d750f21a734b133d8e07e25b98 -address /run/containerd/containerd.sock -containerd-binary /usr/bin/containerd -runtime-root /var/run/docker/runtime-runc -systemd-cgroup
root       1961   1942  0 05:58 ?        00:00:00 nginx: master process nginx -g daemon off;
```

`dockerd`在启动时可以指定`runc`的实现, 使用`docker info`也可以查到`containerd`和`runc`的版本信息.
