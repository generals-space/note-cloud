# CGroup(二)服务的配置文件

参考文章

1. [Linux Cgroup系列（01）：Cgroup概述](https://segmentfault.com/a/1190000006917884)
    - 查看当前进程属于哪些cgroup: `cat /proc/[pid]/cgroup`

## 写在前面

参考文章1中**如何使用cgroup**一节, 给出了`mount`命令挂载指定子系统的示例, 如`mount -t cgroup xxx /sys/fs/cgroup`. 

但是要注意, CentOS 7默认是开启 cgroup 支持的, 无需手动执行这些命令, 很可能会遇到错误

```
$ mount -t cgroup -o cpuset xxx /sys/fs/cgroup/cpuset
mount: none already mounted or /cgroup busy
```

这些子系统其实已经挂载到`/sys/fs/cgroup`目录下了, 使用`mount`命令就可以看到

```
$ mount | grep cgroup
tmpfs on /sys/fs/cgroup type tmpfs (ro,nosuid,nodev,noexec,mode=755)
cgroup on /sys/fs/cgroup/systemd type cgroup (rw,nosuid,nodev,noexec,relatime,xattr,release_agent=/usr/lib/systemd/systemd-cgroups-agent,name=systemd)
cgroup on /sys/fs/cgroup/net_cls,net_prio type cgroup (rw,nosuid,nodev,noexec,relatime,net_prio,net_cls)
...
```

> 在很多使用systemd的系统中，比如ubuntu 16.04，systemd已经帮我们将各个subsystem和cgroup树关联并挂载好了

## mount 子系统到指定目录

按照参考文章1的说法, cgroup 树除了存在于`/sys/fs/cgroup`目录下, 还可以手动挂载到其他路径, 如下

```log
$ mkdir -p /tmp/mycpuset
$ mount -t cgroup -o cpuset mycpuset /tmp/mycpuset
$ mount | grep cpuset
cgroup on /sys/fs/cgroup/cpuset type cgroup (rw,nosuid,nodev,noexec,relatime,cpuset)
mycpuset on /tmp/mycpuset type cgroup (rw,relatime,cpuset)
```

在其中一个目录在做的操作, 另一个目录也能看到, 所以可以判断其实就是同一个目录的两个不同映射.

![](https://gitee.com/generals-space/gitimg/raw/master/4a7fc090cf134d2b76afb52d160551e7.png)

另外, 这两个目录的inode编号其实也是一样的.

![](https://gitee.com/generals-space/gitimg/raw/master/b58a27ee8d168f71a699a207ec0d0bff.png)

清理可以使用`umount`命令.
