# calico-ipip模型分析.1(失败)

参考文章

1. [Calico网络模型](https://www.cnblogs.com/menkeyi/p/11364977.html)
    - calico网络模型, 讲解很详细
    - Pod内部默认网关`169.254.1.1`
2. [Calico 跨子网直连POD网络](https://www.jianshu.com/p/19dca91c71ce)

## 1. 引言

我做了下`ipip`的隧道实验, 但是发现好像和calico构造的网络结构有点出入.

首先, **隧道**理论上应该是点对点的, 下面是我在阿里云上做`ipip`实验的部分命令.

```
ip tunnel add tun_ipip0 mode ipip remote 47.114.45.139 local 172.31.249.7
ip link set tun_ipip0 up
ip addr add 192.168.1.2 peer 192.168.1.1 dev tun_ipip0
ip r add 172.16.144.0/20 dev tun_ipip0
```

`47.114.45.139`为对端主机的公网IP, `172.31.249.7`为当前主机的内网IP(所属网络`172.16.144.0/20`). `192.168.1.1`和`192.168.1.2`为虚构的IP.

这样, 创建的`tun_ipip0`设备的信息如下.

```
56: tun_ipip0@NONE: <POINTOPOINT,NOARP,UP,LOWER_UP> mtu 1480 qdisc noqueue state UNKNOWN group default qlen 1000
    link/ipip 172.31.249.7 peer 47.114.45.139
    inet 192.168.1.2 peer 192.168.1.1 scope global tun_ipip0
       valid_lft forever preferred_lft forever
```

但是`calico`明显创建了一个overlay网络, 各节点的`tunl0`都互相连通, 不是一对一连接. `tunl0`设备的信息上如下.

```
8: tunl0@NONE: <NOARP,UP,LOWER_UP> mtu 1440 qdisc noqueue state UNKNOWN group default qlen 1000
    link/ipip 0.0.0.0 brd 0.0.0.0
    inet 172.16.151.128/32 brd 172.16.151.128 scope global tunl0
       valid_lft forever preferred_lft forever
```

## 2. 实验

### 第1阶段

```bash
ip netns add netns01
ip netns add netns02
ip link add veth12 type veth peer name veth21 
ip link set veth12 netns netns01
ip link set veth21 netns netns02
ip netns exec netns01 ip addr add 172.16.0.1/24 dev veth12
ip netns exec netns02 ip addr add 172.16.0.2/24 dev veth21
ip netns exec netns01 ip link set veth12 up
ip netns exec netns02 ip link set veth21 up
ip netns exec netns02 ip link set lo up
ip netns exec netns01 ip link set lo up
```

```
    +--------------------------------+            +--------------------------------+
    |           +--------+           |            |           +--------+           |
    |           | veth12 |  <------------------------------>  | veth21 |           |
    |           +--------+           |            |           +--------+           |
    |          172.16.0.1/24         |            |          172.16.0.1/24         |
    +--------------------------------+            +--------------------------------+
```

### 第2阶段

```bash
ip netns add netns03
ip netns add netns05
ip link add veth31 type veth peer name veth13
ip link add veth52 type veth peer name veth25
ip link set veth13 netns netns01
ip link set veth31 netns netns03
ip link set veth25 netns netns02
ip link set veth52 netns netns05

ip netns exec netns03 ip addr add 10.10.0.2/24 dev veth31
ip netns exec netns05 ip addr add 10.10.1.2/24 dev veth52
ip netns exec netns01 ip link set veth13 up
ip netns exec netns02 ip link set veth25 up
ip netns exec netns03 ip link set veth31 up
ip netns exec netns05 ip link set veth52 up
ip netns exec netns03 ip link set lo up
ip netns exec netns05 ip link set lo up
```

```
    +--------------+  +--------------+            +--------------+  +--------------+
    |      netns03 |  | netns04      |            |      netns05 |  | netns06      |
    |              |  |              |            |              |  |              |
    | 10.10.0.2/24 |  | 10.10.0.3/24 |            | 10.10.1.2/24 |  | 10.10.1.3/24 |
    |  +--------+  |  |  +--------+  |            |  +--------+  |  |  +--------+  |
    |  | veth31 |  |  |  | veth41 |  |            |  | veth52 |  |  |  | veth62 |  |
    |  +----↑---+  |  |  +---↑----+  |            |  +----↑---+  |  |  +---↑----+  |
    +-------|------+  +------|-------+            +-------|------+  +------|-------+
            |                |                            |                |        
    +-------|----------------|-------+            +-------|----------------|-------+
    |  +----↓---+        +---↓----+  |            |  +----↓---+        +---↓----+  |
    |  | veth13 |        | veth14 |  |            |  | veth25 |        | veth26 |  |
    |  +--------+        +--------+  |            |  +--------+        +--------+  |
    |  netns01                       |            |                       netns02  |
    |           +--------+           |            |           +--------+           |
    |           | veth12 |  <------------------------------>  | veth21 |           |
    |           +--------+           |            |           +--------+           |
    |          172.16.0.1/24         |            |          172.16.0.1/24         |
    +--------------------------------+            +--------------------------------+
```

由于此时`netns01`和`netns02`中的`veth13`和`veth25`没有IP地址, 因此没有生成到`netns03(10.10.10.0/24)`或`netns05(10.10.1.0/24)`的路由. `netns03`中ping不通`netns01`中的`172.16.0.1`, `netns01`中也ping不通`netns03`中的`10.10.0.2`.

按照calico本身的做法, 是使用了一个不存在的IP地址`169.254.1.1`, 并且添加了一条永久的`arp`记录指向veth pair位于宿主机的一端(对应图中的`veth13/veth14`和`veth25/veth26`设备).

这里我们简单一点.

```
ip netns exec netns01 ip r add 10.10.0.2 dev veth13
ip netns exec netns02 ip r add 10.10.1.2 dev veth25
ip netns exec netns03 ip r add 172.16.0.0/24 dev veth31
ip netns exec netns05 ip r add 172.16.0.0/24 dev veth52
```

这样, `netns01`和`netns03`, `netns02`和`netns05`就可以相互ping通了, 不过`netns03`没有办法ping通`netns02`中的`172.16.0.2`, `netns05`也ping不通`netns01`中的`172.16.0.1`.

不过为了接下去的实验能够继续, 还是要将上面的几条先删除.

```
ip netns exec netns01 ip r del 10.10.0.2 dev veth13
ip netns exec netns02 ip r del 10.10.1.2 dev veth25
ip netns exec netns03 ip r del 172.16.0.0/24 dev veth31
ip netns exec netns05 ip r del 172.16.0.0/24 dev veth52
```

### 第3阶段

我们要达到的最终目标是, `netns03`中可以ping通`netns05`中的`10.10.1.2`.

至于上面说的, "`netns03`没有办法ping通`netns02`中的`172.16.0.2`, `netns05`也ping不通`netns01`中的`172.16.0.1`"的情况, 是因为`netns01`与`netns02`是通过veth pair直接相连的, 与实际网络结构不相符, 这个先不用管.

```
ip tunnel add tunl1 mode ipip
ip addr add 10.10.10.1/24 dev tunl1
ip link set tunl1 up
ip link set tunl1 netns netns01
ip netns exec netns01 ip r add 10.10.1.0/24 172.16.0.2 via dev tunl1

ip tunnel add tunl2 mode ipip
ip addr add 10.10.10.1/24 dev tunl2
ip link set tunl2 up
ip link set tunl2 netns netns02
ip netns exec netns02 ip r add 10.10.1.0/24 172.16.0.2 via dev tunl2
```

本来想构建如下结构的网络拓扑的, 结果后来发现上面创建`tunl2`时失败了...

原因好像是linux中`tunnel`设备全局只能存在一个, 一个`tunnel`设备在所有`netns`中都可见, 且真是同一个. 就算`tunl1`是`ipip`, `tunl2`改成`gre`也不行, 显示`add tunnel "tunl0" failed: File exists`.

唉, 放弃了, 换虚拟机来搞.

```
    +--------------+  +--------------+            +--------------+  +--------------+
    |      netns03 |  | netns04      |            |      netns05 |  | netns06      |
    |              |  |              |            |              |  |              |
    | 10.10.0.2/24 |  | 10.10.0.3/24 |            | 10.10.1.2/24 |  | 10.10.1.3/24 |
    |  +--------+  |  |  +--------+  |            |  +--------+  |  |  +--------+  |
    |  | veth31 |  |  |  | veth41 |  |            |  | veth52 |  |  |  | veth62 |  |
    |  +----↑---+  |  |  +---↑----+  |            |  +----↑---+  |  |  +---↑----+  |
    +-------|------+  +------|-------+            +-------|------+  +------|-------+
            |                |                            |                |        
    +-------|----------------|-------+            +-------|----------------|-------+
    |  +----↓---+        +---↓----+  |            |  +----↓---+        +---↓----+  |
    |  | veth13 |        | veth14 |  |            |  | veth25 |        | veth26 |  |
    |  +--------+        +--------+  |            |  +--------+        +--------+  |
    |           +--------+           |            |           +--------+           |
    |           | tunl1  |           |            |           | tunl2  |           |
    |           +--------+           |            |           +--------+           |
    |          10.10.10.1/24         |            |          10.10.1.1/24          |
    |  netns01                       |            |                       netns02  |
    |           +--------+           |            |           +--------+           |
    |           | veth12 |  <------------------------------>  | veth21 |           |
    |           +--------+           |            |           +--------+           |
    |          172.16.0.1/24         |            |          172.16.0.2/24         |
    +--------------------------------+            +--------------------------------+
```
