# calico-ipip模型分析.2

参考文章

1. [Calico网络模型](https://www.cnblogs.com/menkeyi/p/11364977.html)
    - calico网络模型, 讲解很详细
    - Pod内部默认网关`169.254.1.1`
2. [Calico 跨子网直连POD网络](https://www.jianshu.com/p/19dca91c71ce)

## 1. 第1阶段-初始网络结构

A

```bash
ip netns add netns01
ip link add veth1a type veth peer name vetha1
ip link set vetha1 up

ip link set veth1a netns netns01
ip netns exec netns01 ip addr add 10.10.0.2/24 dev veth1a
ip netns exec netns01 ip link set veth1a up
ip netns exec netns01 ip link set lo up
```

B

```
ip netns add netns03
ip link add veth3b type veth peer name vethb3
ip link set vethb3 up
ip link set veth3b netns netns03

ip netns exec netns03 ip link set veth3b up
ip netns exec netns03 ip addr add 10.10.1.2/24 dev veth3b
ip netns exec netns03 ip link set lo up
```

```
    +--------------+  +--------------+            +--------------+  +--------------+
    |      netns01 |  | netns02      |            |      netns03 |  | netns04      |
    |              |  |              |            |              |  |              |
    | 10.10.0.2/24 |  | 10.10.0.3/24 |            | 10.10.1.2/24 |  | 10.10.1.3/24 |
    |  +--------+  |  |  +--------+  |            |  +--------+  |  |  +--------+  |
    |  | veth1a |  |  |  | veth2a |  |            |  | veth3b |  |  |  | veth4b |  |
    |  +----↑---+  |  |  +---↑----+  |            |  +----↑---+  |  |  +---↑----+  |
    +-------|------+  +------|-------+            +-------|------+  +------|-------+
            |                |                            |                |        
    +-------|----------------|-------+            +-------|----------------|-------+
    |  +----↓---+        +---↓----+  |            |  +----↓---+        +---↓----+  |
    |  | vetha1 |        | vetha2 |  |            |  | vethb3 |        | vethb4 |  |
    |  +--------+        +--------+  |            |  +--------+        +--------+  |
    |                                |            |                                |
    |           +--------+           |            |           +--------+           |
    |           |  eth0  |  <------------------------------>  |  eth0  |           |
    |           +--------+           |            |           +--------+           |
    |        172.16.91.201/24        |            |        172.16.91.202/24        |
    +--------------------------------+            +--------------------------------+
```

但是目前主机`A`与其上的`netns`还无法直接通信, 接下来我们按照calico的方式添加路由.

## 2. 第2阶段-路由

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

## 宿主机作为网关
ip r add 10.10.1.0/24 dev ens33 via 172.16.91.202
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

## 宿主机作为网关
ip r add 10.10.0.0/24 dev ens33 via 172.16.91.201
```

此时, `netns`与宿主机, `netns`与`netns`之间都已可以相互通信.

## 3. 第3阶段-tunnel

A

```bash
## 可能会出现这个问题, add tunnel "tunl0" failed: File exists, 但是没关系
ip tunnel add tunl0 mode ipip || echo 'maybe failed but ignore'
ip addr add 10.10.0.1/32 brd 10.10.0.1 dev tunl0
ip link set tunl0 up

ip r del 10.10.1.0/24 dev ens33 via 172.16.91.202
ip r add 10.10.1.0/24 dev tunl0 via 172.16.91.202
```

B

```bash
## 可能会出现这个问题, add tunnel "tunl0" failed: File exists, 但是没关系
ip tunnel add tunl0 mode ipip || echo 'maybe failed but ignore'
ip addr add 10.10.1.1/32 brd 10.10.1.1 dev tunl0
ip link set tunl0 up

ip r del 10.10.0.0/24 dev ens33 via 172.16.91.201
ip r add 10.10.0.0/24 dev tunl0 via 172.16.91.201
```

```
    +--------------+  +--------------+            +--------------+  +--------------+
    |      netns01 |  | netns02      |            |      netns03 |  | netns04      |
    |              |  |              |            |              |  |              |
    | 10.10.0.2/24 |  | 10.10.0.3/24 |            | 10.10.1.2/24 |  | 10.10.1.3/24 |
    |  +--------+  |  |  +--------+  |            |  +--------+  |  |  +--------+  |
    |  | veth1a |  |  |  | veth2a |  |            |  | veth3b |  |  |  | veth4b |  |
    |  +----↑---+  |  |  +---↑----+  |            |  +----↑---+  |  |  +---↑----+  |
    +-------|------+  +------|-------+            +-------|------+  +------|-------+
            |                |                            |                |        
    +-------|----------------|-------+            +-------|----------------|-------+
    |  +----↓---+        +---↓----+  |            |  +----↓---+        +---↓----+  |
    |  | vetha1 |        | vetha2 |  |            |  | vethb3 |        | vethb4 |  |
    |  +--------+        +--------+  |            |  +--------+        +--------+  |
    |                                |            |                                |
    |           +--------+           |            |           +--------+           |
    |           | tunl1  |           |            |           | tunl2  |           |
    |           +--------+           |            |           +--------+           |
    |          10.10.0.1/24          |            |          10.10.1.1/24          |
    |                                |            |                                |
    |           +--------+           |            |           +--------+           |
    |           |  eth0  |  <------------------------------>  |  eth0  |           |
    |           +--------+           |            |           +--------+           |
    |        172.16.91.201/24        |            |        172.16.91.202/24        |
    +--------------------------------+            +--------------------------------+
```

但是在添加`ip r add 10.10.1.0/24 dev tunl0 via 172.16.91.202`这一步会出错: `RTNETLINK answers: Network is unreachable`

calico所在宿主机的路由如下

```
default via 192.168.80.2 dev ens33 proto static metric 100
172.16.36.192/26 via 192.168.80.124 dev tunl0 proto bird onlink
```

其中`172.16.36.192/26`是另一节点上 Pod 的 IP 范围, 该节点的 IP 为 `192.168.80.124`, 默认路由中的`192.168.80.2`为宿主机网络中的网关地址.

我们可以看到, 到对端 Pod 的路由中, `dev`设备为`tunl0`, 但是其`proto`为`bird`. 如果上面的命令中添加了`proto bird`, 就不会再出现`Network is unreachable`的报错了. 不过目前我还不清楚`bird`的生效过程, 应该要看`bird`的源码了, 以后再说吧...

