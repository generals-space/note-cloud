# calico-dual stack双栈

参考文章

1. [官方文档 Enable dual stack](https://docs.projectcalico.org/networking/dual-stack)
2. [官方博客 Enable IPv6 on Kubernetes with Project Calico](https://www.projectcalico.org/enable-ipv6-on-kubernetes-with-project-calico/)
    - 实际操作手册
3. [Configuring calico/node](https://docs.projectcalico.org/reference/node/configuration#environment-variables)
	- calico-node 的环境变量配置

参考文章1中说的步骤很简单, 但也够用了, 参考文章2倒是说得很...啰嗦.

`calico-node`的`DaemonSet`中, `PodCIDR`配置为`fec0:20::/96`.

更新了 ConfigMap, 然后重启各节点上的 calico-node. 不过重启失败了...

```log
$ k logs -n kube-system calico-node-kk2wq
2020-06-06 01:27:07.645 [INFO][8] startup.go 259: Early log level set to info
2020-06-06 01:27:07.646 [INFO][8] startup.go 275: Using NODENAME environment for node name
2020-06-06 01:27:07.646 [INFO][8] startup.go 287: Determined node name: k8s-master-01
2020-06-06 01:27:07.649 [INFO][8] k8s.go 228: Using Calico IPAM
2020-06-06 01:27:07.650 [INFO][8] startup.go 319: Checking datastore connection
2020-06-06 01:27:07.669 [INFO][8] startup.go 343: Datastore connection verified
2020-06-06 01:27:07.669 [INFO][8] startup.go 98: Datastore is ready
2020-06-06 01:27:07.681 [INFO][8] startup.go 621: Using autodetected IPv4 address on interface ens33: 192.168.80.121/24
2020-06-06 01:27:07.682 [WARNING][8] startup.go 617: Unable to auto-detect an IPv6 address: no valid IPv6 addresses found on the host interfaces
2020-06-06 01:27:07.682 [WARNING][8] startup.go 432: Couldn't autodetect an IPv6 address. If auto-detecting, choose a different autodetection method. Otherwise provide an explicit address.
2020-06-06 01:27:07.682 [INFO][8] startup.go 213: Clearing out-of-date IPv4 address from this node IP="192.168.80.121/24"
2020-06-06 01:27:07.682 [INFO][8] startup.go 217: Clearing out-of-date IPv6 address from this node IP=""
2020-06-06 01:27:07.690 [WARNING][8] startup.go 1122: Terminating
Calico node failed to start
```

日志信息不太全, 翻了翻官方文档, 找到了参考文章3, 添加了下`CALICO_STARTUP_LOGLEVEL`环境变量, 写为`DEBUG`, 重启目标 Pod, 这次终于有详细输出了.

```log
2020-06-06 02:33:21.254 [DEBUG][8] interfaces.go 98: Found valid IP address and network CIDR=fe80::20c:29ff:fe81:97/64
2020-06-06 02:33:21.254 [DEBUG][8] filtered.go 41: Check interface Name="ens33"
2020-06-06 02:33:21.254 [DEBUG][8] filtered.go 43: Check address CIDR=fe80::20c:29ff:fe81:97/64
2020-06-06 02:33:21.254 [WARNING][8] startup.go 617: Unable to auto-detect an IPv6 address: no valid IPv6 addresses found on the host interfaces
2020-06-06 02:33:21.254 [WARNING][8] startup.go 432: Couldn't autodetect an IPv6 address. If auto-detecting, choose a different autodetection method. Otherwise provide an explicit address.
```

这里显示calico-node 明明已经检测到网卡上的 IPv6 地址, 但还是没说明出错原因, 后来找到了[相关源码](https://github.com/projectcalico/node/blob/v3.11.0/pkg/startup/autodetection/filtered.go#L40)...

```go
	for _, i := range interfaces {
		log.WithField("Name", i.Name).Debug("Check interface")
		for _, c := range i.Cidrs {
			log.WithField("CIDR", c).Debug("Check address")
			if c.IP.IsGlobalUnicast() {
				return &i, &c, nil
			}
		}
	}
```

其中`c.IP`为`net.IP{}`结构体对象, `IsGlobalUnicast()`会判断此`IP`是否为**全球单播地址**. 这个概念在 IPv4 和 IPv6 中都有, 但是涵义好像有点不太一样.

IPv4 的`192.168.0.0/16`中的局域网地址都属于**全球单播地址**, 但 IPv6 中的`FE80::/10`本地链路地址(对等于`169.254.0.0/16`), 及`FEC0::/10`本地站点地址(对等于`192.168.0.0/16`)却不属于这个范围...

不知道ta们具体的计算方式, 日后再说吧.

我将宿主机(vmware虚拟机)的 IPv6 地址(`fe80::20c:29ff:fe81:97/64`)的前缀, 由`fe80`改成`2001`, 就可以通过这个检查, 然后 calico-node 也启动成功了.

> 由于`ip addr`不支持直接修改地址, 所以本来准备先添加`2001`前缀的地址, 然后把`fe80`的地址删掉的, 但总是删不掉. 后来连 Pod 里也有了`fe80`的地址, 和自己所属的 PodCIDR IPv6 地址一起, 一共两个 IPv6 地址...

之后新创建的 Pod 中, `eth0`网卡上会多出一个 IPv6 的地址, 只不过`kubectl get pod`不能显示出 Pod 的地址来.

但是集群内通过 IPv6 通信还是不行, Pod 无法 ping 通宿主机节点的 IPv6 的地址, Pod 之间也不行.
