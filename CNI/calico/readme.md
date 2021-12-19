# calico

参考文章

1. [Calico网络模型](https://www.cnblogs.com/menkeyi/p/11364977.html)
    - calico网络模型, 讲解很详细
    - Pod内部默认网关`169.254.1.1`
    - 容器跨主机通信的转发流程
    - felix: 路由配置组件; bird: 路由广播组件(BGP Speaker)
    - ipip网络模型的目标.
3. [趣谈网络协议 - 第31讲 | 容器网络之Calico：为高效说出善意的谎言](https://blog.csdn.net/aha_jasper/article/details/105575893)
    - 像是参考文章1的原文
4. [Calico网络方案](https://www.cnblogs.com/netonline/p/9720279.html)
    - 文章末尾给出的calico与flannel各模型网络的对比很值得一看.
5. [calico网络策略](https://yq.aliyun.com/articles/674020)
    - `calico`不同于`flannel`, 不需要为每个node分配子网段, 所以只需要考虑pod的数量;
    - 宿主机上默认初始会创建一个掩码位为26的子网网段, 当pod超过这个值后, 可以从其他可用子网中取值, 并不是固定的.
6. [docker 容器网络方案：calico 网络模型](https://cizixs.com/2017/10/19/docker-calico-network/)
    - calico引入的各组件(`libnetwork-plugin`, `BIRD`, `confd`, `felix`, `etcd`等)及ta们各自的功能.
    - `libnetwork-plugin`是用于与原生docker配合使用的网络插件, 实现的是docker的网络接口. 但kuber集群需要的是CNI插件, 按照readme文档中所说, calico还有一个`cni-plugin`工程, 这才是kuber集群部署需要的.
    - `BIRD`用于宿主机节点间的路由信息的传递, 可以理解为`gossip`, 每创建一个Pod, 就会生成一条到达此Pod的路由.
    - 各组件在容器内部的配置文件目录位置.
    - calico的优点与缺点
7. [DockOne微信分享（二六二）：谐云Calico大规模场景落地实践](http://dockerone.com/article/10382)
8. [Calico IPAM源码解析](https://harmonycloudcaas.yuque.com/docs/share/b4ccb4d5-7b1e-4daf-b0c5-002424073702?#)
    - 泽宇

calico的BGP与flannel的host-gw一样, 是L3的underlay方案.

> 三层通信模型表示每个容器都通过 IP 直接通信, 中间通过路由转发找到对方. --参考文章4. (ARP等协议是不通的, 只有更上层的才可以)

想要让宿主机节点承担起路由器(或者叫网关)的任务, 需要知道每个Pod运行在哪个宿主机节点上, 这样才能为此Pod设置独立的路由.

flannel的host-gw模型是为每个宿主机节点划分一个小子网(16位掩码), 然后为每个子网创建路由. 这样限制了宿主机的数量, 宿主机上的网段空闲的IP就被浪费了, 对集群规模还是有限制的.

calico的BGP则没有这样的限制, 同一个宿主机节点上的Pod也不定属于同一个小子网, ta会针对每个Pod设置独立的路由.

那么flannel是怎么知道每个宿主机节点上的网段, calico又是怎么知道哪个Pod运行在哪台宿主机上呢?

`coreos/flannel`服务会监听宿主机节点的变化, 每新增一个节点, 就为其划分一个新的小子网, 其他已存在的节点上都会添加到新节点Pod子网的路由.

而calico则引入了`BIRD`组件, 通过BGP协议实现路由信息的传播, 类似于gossip, 集群中某一节点的信息更新, 其他节点都会知道. 但是路由信息最终是要在宿主机部署的, 而宿主机上的路由及防火墙操作, 则是通过`fliex`组件完成的.

## IPIP模型

ipip网络模型是为了解决集群中各节点不在同一网段的问题. 一般情况下集群中的节点都在同一局域网, 但也有可能为了做冗余和灾备, 节点位于不同地点的机房, 此时应该是可以通过外网IP实现集群通信的. 但是BGP模型无法在这种类型的网络中通过三层路由完成, 所以ipip出现了.

按照参考文章2末尾所说, calico的ipip与flannel的vxlan/udp一样是overlay的解决方案. 另外也提到, calico的BGP与flannel的host-gw, 都要求集群各节点在相同子网网段中, 也在参考文章1中找到了答案.

## calico网络细节

1. Calico网络中Pod的IP地址掩码位为32, 直接作为单点局域网.
2. Calico网络中Pod的默认路由中, 网关地址为`169.254.1.1`, 在整个网络中这个IP是不存在的.
    - > 当一台机器要访问网关的时候, 首先会通过ARP获得网关的MAC地址, 然后将目标MAC变为网关的MAC, 而网关的IP地址不会在任何网络包头里面出现. 也就是说, 没有人在乎这个地址具体是什么, 只要能找到对应的MAC, 响应ARP就可以了.
3. Calico会在生成的Pod中添加一个永久的arp记录(使用`ip neigh`可以查看), 将`169.254.1.1`指向Pod内部`eth0`的veth pair对端网卡. 从Pod中发出的数据包会将位于宿主机的veth pair设备作为网关, 到达宿主机后, 匹配宿主机上的路由规则. 即, 将宿主机本身视作路由器.
4. 宿主机上拥有到自己本机上所有Pod的路由, 每条路由的目标设备都是位于宿主机一端的veth pair设备.
5. 宿主机上拥有所有到其他宿主机网段的路由, 这些路由的网关IP都是对端宿主机的IP, 而这些路由规则中的`dev`字段根据使用网络模型的不同而不同.
    - IPIP模型, 会创建一个名为`tunl0@NONE`, 类型为`ipip`的设备. 当前节点到其他节点网段的路由规则中, `dev`指向这个设备. 而这个设备会以当前主机被划分子网的网络号为IP. 比如, 当前主机被划分了`172.16.151.128/26`网段, 这个`tunl0`设备的地址为`172.16.151.128/32`, 直接使用网络号, 在集群中可以使用这个IP直接通信.
6. ipip模型下, 跨node的Pod通信会通过`tunl0`设备封装数据包, 而同node中的Pod通信直接走宿主机上的路由即可(`dev`为veth pair中位于宿主机的一端)
