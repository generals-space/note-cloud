# calico网络模型-169.254.1.1的作用.0.环境构建

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

## calico网络细节

calico在Pod内部创建了如下路由

```
default via 169.254.1.1 dev eth0
169.254.1.1 dev eth0 scope link
```

我们知道, `eth0`其实就是veth pair设备位于Pod内的一端, 而calico创建的veth pair设备, 位于宿主机的另一端并没有像flannel那样接入到一个`bridge`网桥, 而是直接留在了宿主机本身.

但是`169.254.1.1`这个地址不存在于宿主机或是任何一个Pod中, 那ta是怎么用的呢?

```

```

参考文章1中说了, "容器里面的默认路由，Calico配置得比较有技巧". 那么这个"技巧"体现在哪呢? 为何要用这种方式呢?

我们知道, 

> 如果全部走三层的路由规则，没必要每台机器都用一个docker0，从而浪费了一个IP地址，而是可以直接用路由转发到veth pair在物理机这一端的网卡。 --参考文章1

...但是总感觉, 为了节省一个bridge设备和一个IP地址, 又是写`arp`规则, 又是重写Pod内部路由的, 太亏了...

## 环境构建

我之前做过通过添加路由规则实现跨主机的docker容器互联的实验(其实就类似flannel的`host-gw`的网络), 但是calico移除了其中的`docker0`网桥, 所以我们这里需要手动构建初始的网络拓扑, 然后再来验证为什么ta的方式"巧妙"(姑且先认为calico节省了一个IP很厉害吧...).

VMware虚拟机环境

- A: 172.16.91.201/24
- B: 172.16.91.202/24

首先将两台主机上的`ip_forward`打开.

```

```

主机`A`上执行如下命令.

```bash
ip netns add netns01
ip link add veth1a type veth peer name vetha1
ip link set vetha1 up

ip link set veth1a netns netns01
ip netns exec netns01 ip addr add 10.10.0.2/24 dev veth1a
ip netns exec netns01 ip link set veth1a up
ip netns exec netns01 ip link set lo up
```

主机`B`上执行如下命令.

```
ip netns add netns03
ip link add veth3b type veth peer name vethb3
ip link set vethb3 up
ip link set veth3b netns netns03

ip netns exec netns03 ip link set veth3b up
ip netns exec netns03 ip addr add 10.10.1.2/24 dev veth3b
ip netns exec netns03 ip link set lo up
```

此时网络拓扑如下, `netns02`与`netns04`为了简单未给出创建步骤, 不过不影响.

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
                    A                                              B                
```

但此时宿主机`A`上ping不通`netns01`, 反过来`netns01`也不ping不通主机`A`, 因为双方都没有到达对方地址的路由. 我们要解决的, 就是如何能让主机A上的`netns`, 能与主机A本身, 与主机B, 及主机B上的`netns`相互通信的问题. 
