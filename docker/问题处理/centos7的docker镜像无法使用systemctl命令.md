# centos7的docker镜像无法使用systemctl命令

参考文章

1. [docker Failed to get D-Bus connection 报错](http://welcomeweb.blog.51cto.com/10487763/1735251)

2. [如何在Docker CentOS容器中使用Systemd](http://www.codesec.net/view/434721.html)

## 问题描述

使用`docker.io/centos:7`的docker官方镜像, 配置`php-fpm`的服务脚本完成. 但是使用`systemctl`命令操作`php-fpm`时出现如下问题. (服务脚本的建立参考[这里](http://www.centoscn.com/CentOS/config/2015/0507/5374.html))

```shell
[root@e13c3d3802d0 /]# systemctl start php-fpm
Failed to get D-Bus connection: Operation not permitted
```

## 原因分析及解决方法

按照参考文章1, 2所说, 这是CentOS:7镜像中的一个bug, 目前无法修复, 只能等到7.2版镜像. 这个bug的原因是因为`dbus-daemon`没能启动. 其实`systemctl`并不是不可以使用. 将你dockerfile的`CMD`或者`ENTRYPOINT`设置为`/usr/sbin/init`即可, 容器会在运行时将`dbus`等服务启动起来. 然后再执行`systemctl`命令即可运行正常.

```shell
docker run --privileged  -e "container=docker"  -v /sys/fs/cgroup:/sys/fs/cgroup -d docker.io/centos:7  /usr/sbin/init
```

> 注意: 上面的命令要完全执行...少一个都会被坑的很惨.

`--privileged`: systemd 依赖于`CAP_SYS_ADMIN capability`. 意味着运行Docker容器需要获得 `privileged`(这不利于一个base image);

`-v`选项: systemd 依赖于访问`cgroups filesystem`; systemd 有很多并不重要的文件存放在一个docker容器中, 如果不删除它们会产生一些错误; 原来以为`-v`选项只是一个共享目录而已, 然而我擅自去掉这个选项后, `systemctl start php-fpm`执行时卡住了, 而且并没有执行成功. 加上这个选项才可以.

另外, 最好加上`-d`选项, 这个是我自己加上的. 交互时启动时竟然提示输入容器密码...而且尝试了很多遍都不对(也不是宿主机的密码)...这真是个奇妙的经历.

```shell
...
CentOS Linux 7 (Core)
Kernel 4.5.5-300.fc24.x86_64 on an x86_64

325606d855b9 login: root
Password:
Login incorrect

325606d855b9 login:
Password:
Login incorrect
...
```