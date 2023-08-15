# kubeadm join无法加入新节点(之前删掉了两个master)

参考文章

1. [Can't join a new node using kubeadm](https://github.com/kubernetes-sigs/kubespray/issues/4117#issuecomment-471232465)
2. [kubeadm join fails if master takes too long to become ready](https://github.com/kubernetes/kubernetes/issues/66768)

场景描述

在用kubeadm创建了3主2从的kuber集群并稳定运行了一段时间后, 由于本地资源没有富余, 就尝试删除了两个主节点, 变成了1主2从.

初时倒也正常, 但是在某次虚拟机休眠并重启后, 发现两个从节点断掉了, 尝试重新用`kubeadm join`加入集群但总是卡住.

```
kubeadm join k8s-server-lb:8443 --token vmiss3.vnakc8bf0gq19ucq     --discovery-token-ca-cert-hash sha256:fa332b427052bfca068567c795f442355fbf37561d54e3fa763992f3cee9a83f
[preflight] Running pre-flight checks
	[WARNING IsDockerSystemdCheck]: detected "cgroupfs" as the Docker cgroup driver. The recommended driver is "systemd". Please follow the guide at https://kubernetes.io/docs/setup/cri/
	[WARNING SystemVerification]: this Docker version is not on the list of validated versions: 19.03.4. Latest validated version: 18.09
^C 卡住
```

然后加入`--v=5`提高日志级别, 看到有如下输出

```
I1130 22:25:16.006053   75102 join.go:433] [preflight] Discovering cluster-info
I1130 22:25:16.006423   75102 token.go:199] [discovery] Trying to connect to API Server "k8s-server-lb:8443"
I1130 22:25:16.007115   75102 token.go:74] [discovery] Created cluster-info discovery client, requesting info from "https://k8s-server-lb:8443"
I1130 22:25:21.018293   75102 token.go:82] [discovery] Failed to request cluster info, will try again: [configmaps "cluster-info" is forbidden: User "system:anonymous" cannot get resource "configmaps" in API group "" in the namespace "kube-public"]
...省略
I1130 22:25:26.019083   75102 token.go:82] [discovery] Failed to request cluster info, will try again: [configmaps "cluster-info" is forbidden: User "system:anonymous" cannot get resource "configmaps" in API group "" in the namespace "kube-public"]
```

然后我找到了参考文章1和2, 说看起来像是master节点没有经过初始化, 导致无法响应这样的join请求.

...就是说我把集群搞挂了呗?

于是我只能把唯一的主节点reset并重新init, 然后再加入两个从节点就可以了.

------

不过理论上通过`kubectl delete node`删除两个主节点应该不会出问题, 我想这可能是之前第一个主节点被我删过一次, 并不是处于正常的运行状态, 增删节点这样的操作是由另外两个主节点完成的, 所以剩下的这个主节点无法处理这样的操作导致了问题. 日后需要重新做一次实验.
