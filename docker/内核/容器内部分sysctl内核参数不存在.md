# 容器内部分sysctl内核参数不存在

尝试在容器内部修改内核参数

```
$ sysctl -w fs.file-max=524287
sysctl: setting key "fs.file-max": Read-only file system
```

容器内部的内核参数可以通过两种方法修改.

一种是在容器启动时通过`--sysctl`选项指定要修改的属性

```
docker run -it --sysctl fs.file-max=40960 --rm generals/centos7 /bin/bash
root@9e850908ddb7:/# sysctl -a | grep fs.file-max
fs.file-max = 40960
```

另一种是在容器启动时设置`--privileged`选项, 之后可以在容器内部使用`sysctl`手动修改

```
docker run -it --privileged --name centos7 generals/centos7 /bin/bash
[root@d25e6ce6cc12 /]# sysctl -w fs.file-max=40960
fs.file-max = 40960
[root@d25e6ce6cc12 /]# sysctl -a | grep fs.file-max
fs.file-max = 40960
```

但是有些内核属性根本就不存在.

场景描述

win10系统, docker桌面版2.0.0.3(build: 8858db3, docker engine: 18.09.2)

```
docker run -it --privileged --name centos7 generals/centos7 /bin/bash
[root@d25e6ce6cc12 /]# sysctl -a | grep backlog
## 无输出
```

如果在cenots7的虚拟机中执行上述命令, 可得到如下输出

```
net.ipv4.tcp_max_syn_backlog = 1024
```

> 修改的内核参数无法保留在容器中, 所以容器重启后要重新执行内核参数的生效操作.
