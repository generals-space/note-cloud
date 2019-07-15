# flannel与calico

参考文章

1. [flannel](https://coreos.com/flannel/docs/latest/kubernetes.html)
2. [Calico](https://docs.projectcalico.org/v3.8/getting-started/kubernetes/installation/calico#installing-with-the-kubernetes-api-datastore50-nodes-or-less)

kubernetes版本: 1.15.0

目前我还不清楚网络插件在kuber集群中是什么角色, 起到了什么作用, 看ta们各自官方的介绍, 貌似是可以脱离kuber独立存在的???

在kuber集群中, 创建[flannel](https://raw.githubusercontent.com/coreos/flannel/master/Documentation/kube-flannel.yml)和[calico](https://docs.projectcalico.org/v3.8/manifests/calico.yaml)都比较方便, 直接使用`kubectl create -f xxx.yaml`就可以了.

不过要注意网络插件与kuber中pod的CIDR要保持一致. 

比如在使用kubeadm init初始化集群时, 可以在命令行中使用`--pod-network-cidr 10.244.0.0/16`选项, 或是在kubeadm配置文件中通过如下格式指定集群内部pod的IP地址范围.

```yaml
networking:
  podSubnet: 10.244.0.0/16
```

那么在创建网络插件时, 这个值也要相同. 在Flannel中, 该配置项在名为`kube-flannel-cfg`的`ConfigMap`中, `net-conf.json`文件的`Network`字段下, 默认为"10.244.0.0/16"; 在Calico中, 此配置可以搜索`CALICO_IPV4POOL_CIDR`, 默认值为"192.168.0.0/16".
