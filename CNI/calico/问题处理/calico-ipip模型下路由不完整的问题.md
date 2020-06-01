# calico-ipip模型下路由不完整的问题

参考文章

1. [bird: Netlink: Network is down](https://github.com/projectcalico/calico/issues/3134)

calico版本: 3.10

ipip网络模型下, 集群中某些节点无法ping通其他节点上的Pod, 当然其上的Pod也和其他节点上的Pod无法通信, 不过有些节点是ok的.

查了查路由, 发现这些有问题的节点上, 到其他节点路由信息不完整. 比如, master节点中有如下路由信息

```
default via 192.168.80.2 dev ens33 proto static metric 100
172.16.36.192/26 via 192.168.80.124 dev ens33 proto bird
172.16.151.128 dev cali720dc2eae25 scope link
blackhole 172.16.151.128/26 proto bird
```

`172.16.151.128/26`是master节点被划分的网段地址, `172.16.151.128`为master节点上某个Pod的IP. 

`172.16.36.192/26`是节点worker01的网段地址, 但是没有worker02的网段路由信息, 因此master只能ping通worker01及其上面的Pod, 但是无法与worker02进行通信.

> 由于`k get node`的状态正常, 所以使用`k exec -it`可以进入到worker02上面的Pod中, 不过`exec`走的是apiserver, 不算与Pod通信, 所以没用.

只有ipip模型有问题, 如果换成bgp模型就一切正常了.

查了下`calico-node`的日志, 发现有`bird: Netlink: Network is down`的日志, 怀疑是这里的问题.

未解决???
