# kubelet无法启动导致无法加入集群

kuber版本: 1.16.2
docker版本: 19.03.4

使用kubeadm建立集群前, 特地按下面的配置修改了docker的启动参数.

```json
{
    "exec-opts": ["native.cgroupdriver=systemd"],
    "log-driver": "json-file",
    "log-opts": {
        "max-size": "100m"
    },
    "storage-driver": "overlay2",
    "storage-opts": [
        "overlay2.override_kernel_check=true"
    ]
}
```

> `dockerd`使用的 cgroup driver 默认为`cgroupfs`.

但在建好集群且正常运行近一个月后的某一天, 节点突然全都失去了联系(之前我把3主2从的集群删了两个主节点, 变成了1主2从), 发现worker节点的kubelet无法启动了, 这也导致无法用`kubeadm join`重新加入集群. 在日志里发现如下错误

> kubelet[60808]: F1130 16:27:12.805286   60808 server.go:271] failed to run Kubelet: failed to create kubelet: misconfiguration: kubelet cgroup driver: "cgroupfs" is different from docker cgroup driver: "systemd"

这明显是说kubelet和docker的cgroup驱动冲突, 但是master节点上的kubelet还运行的好好的, worker节点突然就挂了, 这很让我匪夷所思.

而且kubeadm在进行环境检查的时候明确要求docker使用`systemd`做驱动, 这不是神经病么?

```
[preflight] Running pre-flight checks
	[WARNING IsDockerSystemdCheck]: detected "cgroupfs" as the Docker cgroup driver. The recommended driver is "systemd". Please follow the guide at https://kubernetes.io/docs/setup/cri/
```

这个问题的解决方法很简单, 要么修改kubelet, 要么修改docker, 反正docker默认也是用cgroup做驱动, 改docker的配置就好.

但是这个问题的出现还是让摸不着头脑...

------

对了, 这和`/var/lib/kubelet/config.yaml`文件不存在完全没有关系, 这个文件在`kubeadm join`后会自动生成.

这样kubelet可以启动了, 但是还是无法使用`kubeadm join`集群. 这个问题的最终解决过程另一篇文章, 就是把整个集群都拆了, 唯一的主节点reset再重新init. 而且我又尝试把docker的cgroup驱动改回了systemd, 还成功了. 这说明这个问题根本不是导致`kubeadm join`加入不了集群的原因.

好在问题解决了, 这个问题单独看是没有什么借鉴意义的, 其实并不能影响什么.
