# calico-PodCIDR与宿主机网络重叠引起的问题

参考文章

1. [采坑指南——域名解析问题排查过程](https://cloud.tencent.com/developer/article/1475564)

这只是一个种可能的情况, 没有严格测试过.

本地虚拟机部署的kuber集群, 宿主机网络为`192.168.80.0/24`, 则calico默认的PodCIDR为`192.168.0.0/16`.

在进入到容器内部后, 发现DNS不能工作, 不过可以直接ping通过内网的IP地址.

Pod内`/etc/resolv.conf`的内容如下

```
nameserver 10.96.0.10
search default.svc.cluster.local svc.cluster.local cluster.local
options ndots:5
```

`10.96.0.10`为集群中`kube-dns`的service的地址, 但Pod内也是可以ping通的...

在将PodCIDR网段更换后重新部署, 这个问题就解决了, 因此并未深究.
