# netlink默认路由的操作方法[获取 添加]

参考文章

1. [flannel `v0.10.0` pkg/ip/iface.go -> GetDefaultGatewayIface()](https://github.com/coreos/flannel/blob/v0.10.0/pkg/ip/iface.go)
    - 获取默认路由
2. [kube-ovn `v0.1.0` pkg/daemon/ovs.go -> configureContainerNic()](https://github.com/alauda/kube-ovn/blob/v0.1.0/pkg/daemon/ovs.go)
    - 添加默认路由

## 1. 获取默认路由的方法

```go
	// 获取当前主机的路由列表.
	routes, err := netlink.RouteList(nil, syscall.AF_INET)
	if err != nil {
		return nil, err
	}

	for _, route := range routes {
		// dst 为 nil, 或是为 0.0.0.0/0 的, 其实就是默认路由.
		if route.Dst == nil || route.Dst.String() == "0.0.0.0/0" {
			// 这里的 LinkIndex 学名为接口索引, 使用 `ip link` 等命令的输出中, 
			// 显示结果前面的数字就是各接口的索引值, lo 回环网卡为 1.
			if route.LinkIndex <= 0 {
				return nil, errors.New("Found default route but could not determine interface")
			}
			return net.InterfaceByIndex(route.LinkIndex)
		}
	}

```

代码来源: 

除了判断目标路由的`Dst`成员是否为`nil/0.0.0.0/0`, 也可以通过判断路由对象是否拥有`Gw`成员, 因为只有默认路由才会设置该值.

...多网卡主机应该也可以用此方法判断???

本程序中使用`Gw`作为判断依据.

## 2. 添加默认路由的方法

```go
	_, defaultNet, _ := net.ParseCIDR("0.0.0.0/0")
	err = netlink.RouteAdd(&netlink.Route{
		LinkIndex: containerLink.Attrs().Index,
		Scope:     netlink.SCOPE_UNIVERSE,
		Dst:       defaultNet,
		Gw:        net.ParseIP(gateway),
	})
	if err != nil {
		return fmt.Errorf("config gateway failed %v", err)
	}
```

> 需要传入目标网络的 gateway 对象.

代码来源: 