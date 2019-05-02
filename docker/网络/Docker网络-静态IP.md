# Docker网络-静态IP

参考文章

1. [为Docker容器指定自定义网段的固定IP/静态IP地址](http://blog.csdn.net/gobitan/article/details/51104362)

Docker守护进程启动以后会创建默认网桥`docker0`，其IP网段通常为`172.17.0.0/16`。在启动Container的时候，Docker将从这个网段自动分配一个IP地址作为容器的IP地址。最新版(1.10.3)的Docker内嵌支持在**启动容器**的时候为其指定静态的IP地址。

但在这之前我们要首先创建一个docker网络(不能指定默认的名称为`bridge`的网络). 这个网络与`bridge`, `none`, `host`和`container`同级.

```
$ docker network create --subnet=172.18.0.0/16 mynet
dd8fa1e40dc39844c3990578cde1a21135b538efcb73da4b620f3387c0cd16c6
$ docker network ls
NETWORK ID          NAME                DRIVER              SCOPE
1ab025269f72        bridge              bridge              local               
a4bde17a0747        host                host                local               
dd8fa1e40dc3        mynet               bridge              local               
0ddfd4caff44        none                null                local 
```

`subnet`网段随意指定, 空闲即可(空闲的意思是与另外的docker网段不冲突, 当然, 更不能与物理网段冲突, 不然会出大乱子的);

`mynet`网络名称随意.

------

之后启动容器可以使用`--net`指定目标网络, 并使用`--ip`选项指定目标网络中的某个IP, 通过这种方式启动的容器退出后再使用`start`命令启动, 其IP不变.

```
$ docker run -it --net=mynet --ip=172.18.0.2 daocloud.io/centos:6 /bin/bash
```