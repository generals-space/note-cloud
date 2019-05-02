# Docker网络模式-macvlan

参考文章

1. [CentOS7下docker用原生方法使用宿主机所在网络](http://www.jianshu.com/p/1241ca36687e)

2. [通过MacVLAN实现Docker跨宿主机互联](http://www.10tiao.com/html/357/201704/2247485101/1.html)

`macvlan`同样是docker容器跨主机互联的解决方案, 但它是docker原生支持的, 而不是借助`ovs`, `pipework`等额外的工具. 而且, 通过`macvlan`可以创建与宿主机同网段的网络并且可以和宿主机网络直接通信而不是通过网桥转发.

系统要求

Linux内核版本: v3.9–3.19和4.0+

Docker: 1.12+

宿主机开启数据转发.

```
$ echo 'net.ipv4.ip_forward = 1' >> /etc/sysctl.conf
$ sysctl -p
```

开启网卡混杂模式

```
$ ip link set eno16777736 promisc on
```

首先创建docker网络(与host, bridge, none等网络模式同级)

```
$ docker network create -d macvlan --subnet=172.32.100.0/24 --gateway=172.32.100.2 -o parent=eno16777736 macnet
```

`--subnet`指定与宿主机所在相同网络(没关系, 不会冲突)

`--gateway`指定宿主机网络中的网关

`-o parent`则指定了出口网卡, 还是比较容易理解的.

`macnet`是网络名称, 可随意指定.

```
$ docker network ls
NETWORK ID          NAME                DRIVER              SCOPE
ef2a68c1ca8d        bridge              bridge              local
0166af7b33f8        host                host                local
e312e054c1a0        macnet              macvlan             local
79612cfb6b5a        none                null                local
```

接下来就可以指定目标网络启动容器了.

```
$ docker run -it --net=macnet --ip=172.32.100.248 daocloud.io/centos:6 /bin/bash
```

按照这个方法在多台宿主机上部署, 可以达到跨主机容器互联的目的. 但是, 同一宿主机上的不同容器间可以联通, 不同宿主机上的不同容器也可以联通, 某一宿主机可以与其他宿主机上的容器联通, 但就是宿主机无法与其本身的容器联通...我真是X﹏X

这不科学啊.

关于宿主机无法与其本身容器通信问题暂时还没有头绪, 不知道iptables或者route能不能解决. <???>

------

关于这个macvlan的原理, 我不太懂这个内核模块的作用, 参考文章画了两张图解释, 没看太懂. 这里说一下我的看法.

猜测docker容器启动时分配一个虚拟的mac地址, 然后向本地网关发arp包, 用以说明这个IP由这个虚拟mac占用了, 以后对这个IP的包就发给这个mac地址好了. 类似于arp欺骗, 或者说就是一种虚拟IP的做法.

而宿主机开启了混杂模式, 可以接收到本来不是发送给它自己物理mac地址的数据包, 而当它发现这个数据包的目标mac地址是自己某个docker容器的时候, 将它交由容器处理. 这里的先后顺序不太明白, 也有可能是docker插件了网络栈, 主动将数据包接管, 大概就是这个意思.

**需要说明的一点是, 用docker创建macvlan网络后就无法再用其他网络模式在容器内部访问外网了...!!!**

## FAQ

如果启动的容器时指定的静态IP超过了目标网络范围(subnet指定的值), 会提示错误.

```
$ docker run -it --privileged=true --net=macnet --ip=172.32.1.10 --name=172.32.1.10  --restart=always -v /root/docker_share_dir:/root/docker_share_dir reg01.sky-mobi.com/vpn-reset-service/pptp:1.0.0.0 /bin/bash
docker: Error response from daemon: Invalid address 172.32.1.10: It does not belong to any of this network's subnets.
```