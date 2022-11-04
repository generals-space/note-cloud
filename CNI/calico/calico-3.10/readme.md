参考文章

1. [Configuring calico/node](https://docs.projectcalico.org/v3.10/reference/node/configuration#ip-setting)
2. [k8s网络calico——BGP模式](https://www.cnblogs.com/jinxj/p/9414830.html)
    - calico默认为ipip网络模型, 提到了切换为BGP模型的方法.
3. [calico网络模型中的路由原理](https://segmentfault.com/a/1190000016565044)
    - 关闭ipip网络的方法.
4. [Calico网络方案](https://www.cnblogs.com/netonline/p/9720279.html)
5. [Quickstart for Calico on Kubernetes](https://projectcalico.docs.tigera.io/archive/v3.10/getting-started/kubernetes/)

**CIDR网段的配置在 04-ds.yaml 中, 部署前需要先修改04-ds中的`CALICO_IPV4POOL_CIDR`作为集群内部pod网段.**

`04-ds`类似于`flannel`, 在每个节点上都有运行. `06-deploy`将会部署一个`calico-kube-controllers`容器(master节点).

另外, **calico默认使用`ipip`的网络模型**, 如果要使用`BGP`模型, 则需要修改`04-ds`文件中的`CALICO_IPV4POOL_IPIP`为`Never`(`Off`也可以, 默认为`Always`).

## 网络模型

calico有如下几种网络模型:

ipip: 由`CALICO_IPV4POOL_IPIP`控制, 可选择值: `Always(默认)`, `CrossSubnet`, `Never(或是Off)`.
vxlan: 由`CALICO_IPV4POOL_VXLAN`控制, 可选值: `Always`, `CrossSubnet`, `Never(或是Off, 默认)`.

`CALICO_IPV4POOL_IPIP`和`CALICO_IPV4POOL_VXLAN`是不能共存的. 如果前者的值不为Never, 那么后者就不应再赋值, 同理, 如果后者的值不为Never, 那么前者就不应再赋值.

> 按照参考文章4中所说, `CrossSubnet`应该是`bgp&ipip`的混合方案, **同子网的节点间路由采用bgp, 跨子网的节点间路由采用ipip**.

默认为ipip, 当vxlan被设置为Never, 且ipip也为Never时, calico将采用BGP模型.

ipip模型下, 每个宿主机都会创建一个`tunl0`网络接口. 换用其他模型时此接口不会删掉, 需要手动处理, `ip link del`和`ip tunnel del`都没用, 需要用`modprobe -r ipip`把ipip内核模块移除才可以.
