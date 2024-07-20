# kubeadm双栈升级 - expected at least one IP from each family (v4 or v6) for dual-stack networking

参考文章

1. [Failed to parse subnet when creating IPv4/IPv6 dual stack](https://github.com/kubernetes/kubeadm/issues/1828)
2. [fialed to test IPv6DualStack feature of release version 1.16.0](https://github.com/kubernetes/kubernetes/issues/83006)
    - `serviceSubnet: Invalid value: "10.96.0.0/12,2019:30::/24": couldn't parse subnet`
    - 开启双栈配置后, `networking`字段下, `podSubnet`的解析是正确的, 但是无法解析`serviceSubnet`字段.

- kubeadm(旧): v1.16.2
- kubeadm(新): v1.17.2

```yaml
kind: ClusterConfiguration
kubernetesVersion: v1.16.2
## 用于生成apiserver匹配的证书及kubelet配置文件的请求入口.
controlPlaneEndpoint: "k8s-server-lb:8443"
featureGates:
  ## 开启 IPv6 双栈
  IPv6DualStack: true
networking:
  ## 这两个其实都是默认值
  podSubnet: "10.254.0.0/16"
  serviceSubnet: "10.96.0.0/12"
```

在使用`kubeadm.v1.16.2`执行`upgrade plan`时, 报了如下错误

```
podSubnet: Invalid value: "10.244.10.0/24": expected at least one IP from each family (v4 or v6) for dual-stack networking
```

后来为`podSubnet`和`serviceSubnet`都添加上了 IPv6 的地址, 如下

```yaml
networking:
  podSubnet: 10.254.0.0/16,2019:20::/24
  serviceSubnet: 10.96.0.0/12,2019:30::/24
```

但是还是报错

```
serviceSubnet: Invalid value: "fc00:30::/64,172.30.0.0/16": couldn't parse subnet
```

按照参考文章2中所说, `podSubnet`字段的 IPv6 是可以解析的, 但是`serviceSubnet`不行, 这是`kubeadm`的bug, 需要修复.

最终我并没有先升级`kubeadm`, 而是把`serviceSubnet`的 IPv6 部分删掉了, `upgrade plan`子命令暂时执行成功...

不过后来发现升级 kuber 集群要先升级`kubeadm`, 再执行`upgrade plan`, 我顺序搞反了...

而升级`kubeadm`回来后, 再执行`upgrade plan`, 结果又报

```
serviceSubnet: Invalid value: "10.96.0.0/12": expected at least one IP from each family (v4 or v6) for dual-stack networking
```

这是修复完成了吧...又把`serviceSubnet`的 IPv6 部分加回来了...

