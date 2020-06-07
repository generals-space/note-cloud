# calico-dual stack双栈

`calico-node`中设置的 IPv6 的 PodCIDR 为`fec0:20::/96`.

## 1. Pod 内的网络信息

```console
$ ip a
1: lo: <LOOPBACK,UP,LOWER_UP> mtu 65536 qdisc noqueue state UNKNOWN group default qlen 1000
    link/loopback 00:00:00:00:00:00 brd 00:00:00:00:00:00
    inet 127.0.0.1/8 scope host lo
       valid_lft forever preferred_lft forever
    inet6 ::1/128 scope host
       valid_lft forever preferred_lft forever
4: eth0@if80: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1440 qdisc noqueue state UP group default
    link/ether 76:0d:a8:1a:26:23 brd ff:ff:ff:ff:ff:ff link-netnsid 0
    inet 172.16.151.160/32 scope global eth0
       valid_lft forever preferred_lft forever
    inet6 fec0:20::4916:a300/128 scope site
       valid_lft forever preferred_lft forever
    inet6 fe80::740d:a8ff:fe1a:2623/64 scope link
       valid_lft forever preferred_lft forever
$ ip r
default via 169.254.1.1 dev eth0
169.254.1.1 dev eth0 scope link
$ ip -6 r
fe80::/64 dev eth0 proto kernel metric 256 pref medium
fec0:20::4916:a300 dev eth0 proto kernel metric 256 pref medium
default via fe80::ecee:eeff:feee:eeee dev eth0 metric 1024 pref medium
```

`eth0`网卡上的 IPv6 地址`fec0:20::4916:a300/128`, 掩码位128, 和 IPv4 中的32一样.

在未开启 IPv6 双栈前, `calico`在为 Pod 创建的 veth pair 位于宿主机这一端的设备赋予了`169.254.1.1/16`. calico 会自动为 Pod 添加`169.254.1.1`到`ee:ee:ee:ee:ee:ee`的arp记录. 

但开了双栈后, 貌似并不会自动添加 IPv6 地址与 mac 地址的 arp 记录, veth pair 位于宿主机这一端的地址也就只有 IPv6 的了.

## 2. 宿主机上的网络信息.

```console
$ ip a
1: lo: <LOOPBACK,UP,LOWER_UP> mtu 65536 qdisc noqueue state UNKNOWN group default qlen 1000
    link/loopback 00:00:00:00:00:00 brd 00:00:00:00:00:00
    inet 127.0.0.1/8 scope host lo
       valid_lft forever preferred_lft forever
    inet6 ::1/128 scope host
       valid_lft forever preferred_lft forever
2: ens33: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1500 qdisc pfifo_fast state UP group default qlen 1000
    link/ether 00:0c:29:95:16:5a brd ff:ff:ff:ff:ff:ff
    inet 192.168.80.125/24 brd 192.168.80.255 scope global noprefixroute ens33
       valid_lft forever preferred_lft forever
    inet6 2001::20c:29ff:fe95:165a/64 scope global
       valid_lft forever preferred_lft forever
    inet6 fe80::20c:29ff:fe95:165a/64 scope link noprefixroute
       valid_lft forever preferred_lft forever
22: calic0b145da649@if4: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1440 qdisc noqueue state UP group default
    link/ether ee:ee:ee:ee:ee:ee brd ff:ff:ff:ff:ff:ff link-netnsid 0
    inet6 fe80::ecee:eeff:feee:eeee/64 scope link
       valid_lft forever preferred_lft forever

$ ip -6 r
2001::/64 dev ens33 proto kernel metric 256 pref medium
fe80::/64 dev ens33 proto kernel metric 100 pref medium
blackhole fec0:20::1925:2800/122 dev lo proto bird metric 1024 error -22 pref medium
fec0:20::4916:a300/122 via 2001::20c:29ff:fe81:97 dev ens33 proto bird metric 1024 pref medium
fec0:20::4bce:9a40/122 via 2001::20c:29ff:fe05:62ec dev ens33 proto bird metric 1024 pref medium
```

宿主机的 IPv6 路由部分有添加到其他宿主机上 Pod 的路由, 但貌似没有到自己本身网段 Pod 的路由...

我尝试手动添加了宿主机到 Pod 的路由, 但没过多久就被移除了...???

在路由有效期间, Pod 与 Pod 之间, Pod 与宿主机节点之间还是可以相互通信的.

> 注意需要开启 IPv6 的路由转发`net.ipv6.conf.all.forwarding`, 与 IPv4 的`net.ipv4.ip_forward`不是同一个参数.