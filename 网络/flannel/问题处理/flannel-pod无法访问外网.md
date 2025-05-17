# kuber-pod无法访问外网(flannel插件)

参考文章

1. [iptables snat规则缺失，kubernetes集群问题node上所有容器无法ping通外网](https://www.codercto.com/a/59991.html)
2. [kubernetes-在pod里面的容器不能ping外部ip](https://blog.csdn.net/kozazyh/article/details/80595782)

kube版本: 1.16.2
flannel: 0.11
网络模式: vxlan

## 问题描述

kuberntes安装完成, 部署flannel插件后, 启动一个alpine的DaemonSet进行测试. 

进入到pod中后发现无法ping通外网, 但是可以pod之间可以ping通, pod与宿主机也可以ping通.

## 原因分析

> 20250427更新: 有可能是因为 flannel 的 configMap 中配置的 pod 网段与 apiserver 中的不一致...

按照参考文章1和2中介绍, 是iptables缺少了一条.

在修改之前的 nat表 -> POSTROUTING链 如下

```log
Chain POSTROUTING (policy ACCEPT)
num  target     prot opt source               destination
1    KUBE-POSTROUTING  all  --  0.0.0.0/0            0.0.0.0/0            /* kubernetes postrouting rules */
2    MASQUERADE  all  --  172.17.0.0/16        0.0.0.0/0
```

## 解决方法

按照参考文章2中所说的

```
iptables -t nat -I POSTROUTING -s 10.254.0.0/16 -j MASQUERADE
```

> 注意: 参考文章1中的是`-A`追加, 但是这样没用, 只有`-I`可以.

修改后的结果为

```
Chain POSTROUTING (policy ACCEPT)
target     prot opt source               destination
MASQUERADE  all  --  10.254.0.0/16        0.0.0.0/0
KUBE-POSTROUTING  all  --  0.0.0.0/0            0.0.0.0/0            /* kubernetes postrouting rules */
MASQUERADE  all  --  172.17.0.0/16        0.0.0.0/0
```
