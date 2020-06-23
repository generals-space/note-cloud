# coredns无法解析的可能原因

1. 节点间的网络已断.

有次在 kuber 集群中部署 es 集群, 发现各 es 实例都找不到其他节点, 最初以为是 es 配置问题. 

后来在排查果发现连 kube-dns 服务也无法连接, 进而发现无法 ping 通外网 IP(不过外网IP还是可以的). 此时目标锁定在 coredns 解析的问题上, 但是 Pod 内部的`/etc/resolv.conf`并无问题. 

正要排查 coredns 的解析日志, 检查是否可能出现错误时, 发现从 master 已经 ping 不通其他 worker 节点上的 Pod 了, 不知道是否是因为频繁挂起虚拟机导致的问题.

重启 worker 节点, 然后将这些节点上的 Pod 删除重建, 但都无济于事, 所以只好 reset 集群了...
