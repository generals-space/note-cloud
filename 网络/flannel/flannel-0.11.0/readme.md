参考文章

1. [Quickstart for Calico on Kubernetes](https://docs.projectcalico.org/v3.10/getting-started/kubernetes/)
    - 用于kuber集群的部署文件

flannel中有一种网络模型是host-gw, 其实就是把宿主机当作网关的方式, 与calico的bgp原理基本相同, 都是三层的解决方案.

在使用`udp`模型时, 需要修改`02-cm.yaml`文件中的`Type`字段, 同时将`03-ds.yaml`中的`securityContext`设置为`privileged: true`, 否则Pod可能启动失败, 报如下错误.

```console
$ k logs -f kube-flannel-ds-amd64-mng7n
...
I0424 10:54:03.024305       1 main.go:386] Found network config - Backend type: udp
E0424 10:54:03.024484       1 main.go:289] Error registering network: failed to open TUN device: open /dev/net/tun: no such file or directory
I0424 10:54:03.024553       1 main.go:366] Stopping shutdownHandler...
```
