# calico网络模型-169.254.1.1的作用.1.模拟

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

A

```console
$ ip r add 10.10.0.2 dev vetha1
$ ip netns exec netns01 ip r add 172.16.91.0/24 dev veth1a
$ ip r add 10.10.1.0/24 dev ens33 via 172.16.91.202
## 现在 netns 中的路由如下
$ ip netns exec netns01 ip r 
10.10.0.0/24 dev veth1a proto kernel scope link src 10.10.0.2 
172.16.91.0/24 dev veth1a scope link 
```

B

```console
$ ip r add 10.10.1.2 dev vethb3
$ ip netns exec netns03 ip r add 172.16.91.0/24 dev veth3b
$ ip r add 10.10.0.0/24 dev ens33 via 172.16.91.201
## 现在 netns 中的路由如下
$ ip netns exec netns03 ip r
10.10.1.0/24 dev veth3b proto kernel scope link src 10.10.1.2 
172.16.91.0/24 dev veth3b scope link 
```

现在, 主机`A`与其上的`netns01`已经能够互相ping通了, 但是`netns01`还是没办法`ping`通过宿主机网关`172.16.91.2`, 也没有办法ping通主机`B`, 就更别说主机`B`上的`netns03`了.

另外, 主机`A`也没有办法ping通主机`B`上的`netns03`, 虽然`A`已经设置了到对方的路由, 但是这条通信链路从`A.ens33 -> B.ens33 -> B.vethb3 -> B.veth3b`还是通的, 但是`B.veth3b`没有回应, 因为`netns03`中到宿主机网络`172.16.91.0/16`的路径还是断的.

这下我们发现, 问题就出在`netns01`没办法以其宿主机`A`为网关, 将请求转发到其他宿主机节点. 因为我们在`netns01`中设置了到宿主机网络的路由`ip r add 172.16.91.0/24 dev veth1a`, 怎么还可能再将主机`A`当作网关呢?

```console
$ ip netns exec netns01 ip r add 172.16.91.0/24 via 172.16.91.201 dev veth1a
RTNETLINK answers: File exists
```

想一想这个过程, 我们要在`netns01`中先ping通过主机`A`, 就要添加到主机`A`所在网络的路由`172.16.91.0/24 dev veth1a`, 添加完这条路由后, 还要再与其他宿主机节点通信, 必然需要以`A`作网关, 但是这已经与前面的路由重复, 因而添加失败了...

------

OK, 如果单独把`netns01`到主机的路由拿出来呢? 将到`172.16.91.201`和`172.16.91.202`单独写一条路由, 然后再写`172.16.91.0/24`, 网关就分别设置为`172.16.91.201`和`172.16.91.202`.

A

```console
$ ip netns exec netns01 ip r del 172.16.91.0/24 dev veth1a
$ ip netns exec netns01 ip r add 172.16.91.201 dev veth1a
$ ip netns exec netns01 ip r add 172.16.91.0/24 dev veth1a via 172.16.91.201
## 现在 netns 中的路由如下
$ ip netns exec netns01 ip r
10.10.0.0/24 dev veth1a proto kernel scope link src 10.10.0.2 
172.16.91.0/24 via 172.16.91.201 dev veth1a 
172.16.91.201 dev veth1a scope link 
```

B

```console
$ ip netns exec netns03 ip r del 172.16.91.0/24 dev veth3b
$ ip netns exec netns03 ip r add 172.16.91.202 dev veth3b
$ ip netns exec netns03 ip r add 172.16.91.0/24 dev veth3b via 172.16.91.202
## 现在 netns 中的路由如下
$ ip netns exec netns03 ip r
10.10.1.0/24 dev veth3b proto kernel scope link src 10.10.1.2 
172.16.91.0/24 via 172.16.91.202 dev veth3b 
172.16.91.202 dev veth3b scope link 
```

这下可以了, 所有链路都通了. 

接下来验证同主机上不同`netns`间的通信.

### netns访问公网

如果想要在`netns`中访问公网, 还需要添加一条默认路由, 仍用宿主机做网关.

A

```
ip netns exec netns01 ip r add default dev veth1a via 172.16.91.201
```

B

```
ip netns exec netns03 ip r add default dev veth3b via 172.16.91.202
```

### 同主机不同netns间通信

其实 netns 中有了上面的, 使用宿主机作为网关的默认路由, 就不再需要执行这里的步骤了.

先在主机`A`上再创建`netns02`

```bash
ip netns add netns02
ip link add veth2a type veth peer name vetha2
ip link set vetha2 up

ip link set veth2a netns netns02
ip netns exec netns02 ip addr add 10.10.0.3/24 dev veth2a
ip netns exec netns02 ip link set veth2a up
ip netns exec netns02 ip link set lo up
```

然后为`netns01`和`netns02`添加到对方的, 以主机`A`为网关的路由.

```bash
ip r add 10.10.0.2 dev vetha1
ip r add 10.10.0.3 dev vetha2

ip netns exec netns01 ip r del 10.10.0.0/24 dev veth1a 
ip netns exec netns02 ip r del 10.10.0.0/24 dev veth2a 

## 主要是这里
ip netns exec netns01 ip r add 172.16.91.201 dev veth1a
ip netns exec netns01 ip r add 10.10.0.0/24 dev veth1a via 172.16.91.201
ip netns exec netns02 ip r add 172.16.91.201 dev veth2a
ip netns exec netns02 ip r add 10.10.0.0/24 dev veth2a via 172.16.91.201
```
