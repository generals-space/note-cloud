# calico网络模型-169.254.1.1的作用.2.calico的做法

参考文章

1. [Calico网络模型](https://www.cnblogs.com/menkeyi/p/11364977.html)
    - calico网络模型, 讲解很详细
    - Pod内部默认网关`169.254.1.1`
2. [（原创）RFC3927研究笔记（RFC3927，local link addr，LLA）](https://www.cnblogs.com/liu_xf/archive/2012/05/26/2519345.html)
    - LLA(Link Local Address), 链路本地地址, 是设备在本地网络中通讯时用的地址. 网段为`169.254.0.0/16`
    - LLA是本地链路的地址, 是在本地网络通讯的, **不通过路由器转发**, 因此网关为0.0.0.0.
    - LLA在分配时的具体流程: PROBING -> ANNOUNCING -> BOUND
3. [169.254.0.0/16这段地址用途](https://blog.csdn.net/onwer3/article/details/47339469)
    - `169.254.0.0/16`属于B类网络
4. [IPv4地址分类（A类 B类 C类 D类 E类）](https://blog.csdn.net/D_R_L_T/article/details/96606543)

现在我们看看calico是怎么做的, 首先记得将前面实验中弄乱的实验环境恢复.

A

```bash
## 添加宿主机到 netns 的路由
ip r add 10.10.0.2 dev vetha1

## 修改 netns 内部的路由
ip netns exec netns01 ip r del 10.10.0.0/24 dev veth1a 
ip netns exec netns01 ip r add 169.254.1.1 dev veth1a scope link
ip netns exec netns01 ip r add default via 169.254.1.1 dev veth1a

## 修改 vetha1 的 mac 地址
ip link set addr ee:ee:ee:ee:ee:ee dev vetha1

## 在 netns 中添加永久 arp 记录, 映射 169.254.1.1 的 mac 地址
## calico 中这条记录的标记为`stale`而非`permanent`, 但是`stale`的话, 这个条会经常发生变动,
## 貌似是 calico 服务会定时维护这个记录, 才能保持ta的不变.
ip netns exec netns01 ip neigh add 169.254.1.1 dev veth1a lladdr ee:ee:ee:ee:ee:ee nud permanent

```

B

```bash
## 添加宿主机到 netns 的路由
ip r add 10.10.1.2 dev vethb3

## 修改 netns 内部的路由
ip netns exec netns03 ip r del 10.10.1.0/24 dev veth3b 
ip netns exec netns03 ip r add 169.254.1.1 dev veth3b scope link
ip netns exec netns03 ip r add default via 169.254.1.1 dev veth3b

## 修改 vethb3 的 mac 地址
ip link set addr ee:ee:ee:ee:ee:ee dev vethb3

## 在 netns 中添加永久 arp 记录, 映射 169.254.1.1 的 mac 地址
## calico 中这条记录的标记为`stale`而非`permanent`, 但是`stale`的话, 这个条会经常发生变动,
## 貌似是 calico 服务会定时维护这个记录, 才能保持ta的不变.
ip netns exec netns03 ip neigh add 169.254.1.1 dev veth3b lladdr ee:ee:ee:ee:ee:ee nud permanent

```

现在`netns`已经可以和宿主机网络直接通信了, 不过还不能实现跨主机的`netns`互通, 因为宿主机上还没有添加到其他主机`netns`的路由.

A

```
ip r add 10.10.1.0/24 dev ens33 via 172.16.91.202
```

B

```
ip r add 10.10.0.0/24 dev ens33 via 172.16.91.201
```

现在可以了.

