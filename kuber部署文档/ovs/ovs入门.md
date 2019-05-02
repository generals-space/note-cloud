参考文章

1. [Open vSwitch2.3.0版本安装部署及基本操作](http://www.sdnlab.com/3166.html)

2. [官方安装手册](http://docs.openvswitch.org/en/latest/intro/install/general/)

3. [网络虚拟化技术（一）: Linux网络虚拟化](https://blog.kghost.info/2013/03/01/linux-network-emulator/)

`ovs-vsctl add-br 网桥名`

`brctl addbr 网桥名`

重启网络服务, 创建的网桥和veth对不会消失, 只有重启系统才能恢复.

## FAQ

### 1. connection attempt failed (Address family not supported by protocol

```
$ ovs-vsctl add-br br0
ovs-vsctl: Error detected while setting up 'br0'.  See ovs-vswitchd log for details.
$ ovs-vswitchd log
2017-06-08T06:43:28Z|00001|ovs_numa|INFO|Discovered 4 CPU cores on NUMA node 0
2017-06-08T06:43:28Z|00002|ovs_numa|INFO|Discovered 1 NUMA nodes and 4 CPU cores
2017-06-08T06:43:28Z|00003|reconnect|INFO|log: connecting...
2017-06-08T06:43:28Z|00004|reconnect|INFO|log: connection attempt failed (Address family not supported by protocol)
2017-06-08T06:43:28Z|00005|reconnect|INFO|log: waiting 1 seconds before reconnect
2017-06-08T06:43:29Z|00006|reconnect|INFO|log: connecting...
2017-06-08T06:43:29Z|00007|reconnect|INFO|log: connection attempt failed (Address family not supported by protocol)
2017-06-08T06:43:29Z|00008|reconnect|INFO|log: waiting 2 seconds before reconnect
2017-06-08T06:43:30Z|00009|fatal_signal|WARN|terminating with signal 2 (Interrupt)
```

情景描述: 

编译安装完成ovs后一切正常, 重启了一下, 再次操作时就成这样了...

尝试过再次重启, 也使用`ovs-vsctl emer-reset`, 都不管用.

解决方法:

记得加载`openvswitch`内核模块(′⌒\`), 然后重启`ovsdb-server`与`ovs-vswitchd`服务.

### 2. 

参考文章

[docker -实际应用常见问题总结FAQ](http://blog.163.com/weixia_1985/blog/static/96304797201649105218139/)

```
$ ifup kbr0
Bringing up interface kbr0:  Error: Connection activation failed: Failed to determine connection's virtual interface name
```

情景描述

使用`brctl addbr kbr0`创建网桥后, 写入`/etc/sysconfig/network-scripts/ifcfg-kbr0`文件, 重启network服务失败. 查看日志时报上述错误.

解决方法:

关闭NetworkManager

```
$ systemctl stop NetworkManager 
chkconfig NetworkManager off
$ systemctl disable NetworkManager
```