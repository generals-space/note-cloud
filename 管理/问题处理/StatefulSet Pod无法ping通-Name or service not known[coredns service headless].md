# StatefulSet Pod无法ping通-Name or service not known[coredns]

## 问题描述

在集群default命名空间中创建了4个statefulset, 每个 sts 的副本数为 1, 但是这几个 Pod 之间无法互相 ping 通, 按照如下格式写全域名也不行.

`pod名称.service名称.ns名称.svc.cluster.local`

不只是这几个 pod 之间无法互相ping通, 在宿主机上也无法 ping 通这些 pod, ta们只能自己 ping 自己.

coredns功能没有问题, 因为其他ns中的 sts pod 是可以ping通的.

## 解决方法

1. sts需要指定 serviceName (否则会报错);
2. 对应的 service 资源必须创建;
    - sts 只要指定了 serviceName 就可以创建成功, 但是不要求对应的 service 已经存在;
3. service 必须显示设置为`clusterIP: None`;
4. sts pod 中必须要保证对应 service 中声明的 port 端口已启动;
    - 如无端口映射只用 IP 进行通信, 可以把 service 的`spec.port:`部分删掉;
