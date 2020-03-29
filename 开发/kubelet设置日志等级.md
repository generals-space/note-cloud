# kubelet设置日志等级

参考文章

1. [kubelet fails to get cgroup stats for docker and kubelet services](https://stackoverflow.com/questions/46726216/kubelet-fails-to-get-cgroup-stats-for-docker-and-kubelet-services)

- kuber部署方式: kubeadm
- kubelet: 1.16.2

`kubelet -h`输出中有`--v=5`的方法可设置日志级别, 而且这个选项不能写在`--config=/var/lib/kubelet/config.yaml`配置文件中.

我曾经尝试在`/usr/lib/systemd/system/kubelet.service`为kubelet的启动命令直接添加`-v/--v`, 也无效.

```ini
[Service]
ExecStart=/usr/bin/kubelet --v=5
Restart=always
StartLimitInterval=0
RestartSec=10
```

当时在网上到处搜索**kubelet 日志等级**, **kubelet log level**都没找到结果, 最后却在一次意外中发现了, 就是参考文章中的报错. 报错的解决方法很简单, 照做就行了, 这里主要讲一下kubelet的配置方法.

`kubelet.service`文件中`ExecStart`只有`ExecStart=/usr/bin/kubelet`, 但实际还存在一个`kubelet.service.d`目录, 该目录下初始存在一个`10-kubeadm.conf`文件, 内容如下.

```ini
# Note: This dropin only works with kubeadm and kubelet v1.11+
[Service]
Environment="KUBELET_KUBECONFIG_ARGS=--bootstrap-kubeconfig=/etc/kubernetes/bootstrap-kubelet.conf --kubeconfig=/etc/kubernetes/kubelet.conf"
Environment="KUBELET_CONFIG_ARGS=--config=/var/lib/kubelet/config.yaml"
# This is a file that "kubeadm init" and "kubeadm join" generates at runtime, populating the KUBELET_KUBEADM_ARGS variable dynamically
EnvironmentFile=-/var/lib/kubelet/kubeadm-flags.env
# This is a file that the user can use for overrides of the kubelet args as a last resort. Preferably, the user should use
# the .NodeRegistration.KubeletExtraArgs object in the configuration files instead. KUBELET_EXTRA_ARGS should be sourced from this file.
EnvironmentFile=-/etc/sysconfig/kubelet
ExecStart=
ExecStart=/usr/bin/kubelet $KUBELET_KUBECONFIG_ARGS $KUBELET_CONFIG_ARGS $KUBELET_KUBEADM_ARGS $KUBELET_EXTRA_ARGS
```

这下好了, 原来实际的`ExecStart`在这里, 选项参数也是在这里设置的.

我们不修改原来的`Environment`, 只在`EnvironmentFile=-/etc/sysconfig/kubelet`环境文件中做修改, 其默认内容如下

```
KUBELET_EXTRA_ARGS=
```

修改为如下

```
KUBELET_EXTRA_ARGS= --v=5
```

然后`systemctl restart kubelet`重启kubelet组件, 使用ps查看会发现, 多了`--v=5`选项.

```
$ ps -ef | grep kubelet
root      31610      1  3 12:01 ?        00:00:00 /usr/bin/kubelet --bootstrap-kubeconfig=/etc/kubernetes/bootstrap-kubelet.conf --kubeconfig=/etc/kubernetes/kubelet.conf --config=/var/lib/kubelet/config.yaml --cgroup-driver=systemd --network-plugin=cni --pod-infra-container-image=registry.cn-hangzhou.aliyuncs.com/google_containers/pause:3.1 --v=5
```

不过也看不出来有没有生效, `/var/log/message`中的输出都是`I0328`, 因为代码里用的就是`klog.V(3).Info()`什么的, 不像用`Debug`来标识.

以后再说吧.
