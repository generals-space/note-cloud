# calico-ipip模型分析.2

参考文章

1. [Calico网络模型](https://www.cnblogs.com/menkeyi/p/11364977.html)
    - calico网络模型, 讲解很详细
    - Pod内部默认网关`169.254.1.1`
2. [Calico 跨子网直连POD网络](https://www.jianshu.com/p/19dca91c71ce)

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

由于宿主机A和B上此时`netns01`和`netns03`中的`vetha1`和`vethb3`没有IP地址, 因此没有生成到`netns01(10.10.10.0/24)`或`netns03(10.10.1.0/24)`的路由. `netns01`中ping不通宿主机`A`中的`172.16.91.201`, 宿主机`A`也ping不通`netns01`中的`10.10.0.2`.

按照calico本身的做法, 是使用了一个不存在的IP地址`169.254.1.1`, 并且添加了一条永久的`arp`记录指向veth pair位于宿主机的一端(对应图中的`vetha1/vetha2`和`vethb3/vethb4`设备).

这里我们简单一点.

```
ip r add 10.10.0.2 dev vetha1
ip netns exec netns01 ip r add 172.16.91.0/24 dev veth1a
```

```
ip r add 10.10.1.2 dev vethb3
ip netns exec netns03 ip r add 172.16.91.0/24 dev veth3b
```

这样, `netns01`和`netns01`, `netns02`和`netns03`就可以相互ping通了, 不过`netns01`没有办法ping通`netns02`中的`172.16.0.2`, `netns03`也ping不通`netns01`中的`172.16.91.201`.

不过为了接下去的实验能够继续, 还是要将上面的几条先删除.

### 第3阶段

我们要达到的最终目标是, `netns01`中可以ping通`netns03`中的`10.10.1.2`.

至于上面说的, "`netns01`没有办法ping通`netns02`中的`172.16.0.2`, `netns03`也ping不通`netns01`中的`172.16.91.201`"的情况, 是因为`netns01`与`netns02`是通过veth pair直接相连的, 与实际网络结构不相符, 这个先不用管.

A

```bash
ip tunnel add tunl0 mode ipip || echo 'maybe failed but ignore'
ip addr add 10.10.1.1/32 dev tunl0
ip link set tunl0 up
ip r add 10.10.1.0/24 via 172.16.91.202 dev ens33
```

B

```bash
ip tunnel add tunl0 mode ipip || echo 'maybe failed but ignore'
ip addr add 10.10.2.1/32 dev tunl0
ip link set tunl0 up
ip r add 10.10.0.0/24 via 172.16.91.201 dev ens33
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
    |           +--------+           |            |           +--------+           |
    |           | tunl1  |           |            |           | tunl2  |           |
    |           +--------+           |            |           +--------+           |
    |          10.10.10.1/24         |            |          10.10.1.1/24          |
    |  netns01                       |            |                       netns02  |
    |           +--------+           |            |           +--------+           |
    |           | veth12 |  <------------------------------>  | veth21 |           |
    |           +--------+           |            |           +--------+           |
    |          172.16.91.201/24         |            |          172.16.0.2/24         |
    +--------------------------------+            +--------------------------------+
```
