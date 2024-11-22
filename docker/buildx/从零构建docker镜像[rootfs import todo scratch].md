# 从零构建docker镜像[rootfs]

centos: 7
docker: 19.03.5

docker镜像解压开的rootfs只是一组常规目录, 由镜像启动为容器时, 只是在完成namespace隔离后, 使用类似 chroot 的方式将 root 路径指向了镜像的 rootfs 而已.

现在尝试一下构建简单的 rootfs, 将会借助 yum 工具.

## 1. 

首先创建一个 rootfs 目录, 用于存放{bin,etc,lib,usr,var}等目录.

然后使用 yum 下载核心工具库, 下载的过程中, yum 工具会在 rootfs 下自动按照 lib, var 的组织结构存放.so共享库, 各种依赖, 缓存等.

```bash
mkdir /root/rootfs
yum install -y --installroot=/root/rootfs coreutils os-release
```

这就是一个简单的镜像了.

```console
$ ls /root/rootfs
bin  boot  dev  etc  home  lib  lib64  media  mnt  opt  proc  root  run  sbin  srv  sys  tmp  usr  var
```

打包, 然后使用 docker 导入.

```bash
cd /root/rootfs
## 注意打包路径, 得到的 tar.gz 包解压出来不是 rootfs 目录, 而是上面的{bin, boot, dev ...}等一大堆目录
tar -czf ../centos.tar.gz ./
cd /root/
docker import centos.tar.gz centos:7-rootfs
```

但是这样得到的镜像在启动时有点问题

```log
$ docker run -it --name centos centos:7-rootfs bash
bash-4.2# ls
bin  boot  dev  etc  home  lib  lib64  media  mnt  opt  proc  root  run  sbin  srv  sys  tmp  usr  var
```
