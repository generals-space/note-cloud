# scheduler启动时端口冲突 - --port参数不生效[10251]

参考文章

1. [kube-scheduler ignoring settings and listening on insecure port](https://github.com/kubernetes-sigs/kubespray/issues/7739)

kube: 1.17.2

## 场景描述

单节点 k8s 集群, 在主节点上通过vscode远程运行 kube-scheduler 工程.

由于 kubeadm 创建的 k8s 集群中, kube-scheduler 已经占用了2个原生端口(10251与10259), 所以需要为我们的工程指定其他端口防止冲突.

```console
$ netstat -nlpt | grep schedule
tcp     0    0 127.0.0.1:10259    0.0.0.0:*    LISTEN    76970/kube-schedule
tcp6    0    0 :::10251           :::*         LISTEN    76970/kube-schedule
```

在 vscode 中, 我配置了如下字段

```json
"args": [
    "--address=0.0.0.0",
    "--port=11251",         // 防止与原生 scheduler 端口冲突
    "--bind-address=0.0.0.0",
    "--secure-port=11259", // 防止与原生 scheduler 端口冲突
    "--kubeconfig=/etc/kubernetes/scheduler.conf",
    "--config=${workspaceFolder}/.vscode/kube-scheduler/scheduler-config.yaml",
    "--policy-config-file=${workspaceFolder}/.vscode/kube-scheduler/scheduler-policy.json",
]
```

但在启动时, 仍然显示监听10251, 然后报错了.

```log
I0929 09:34:43.954974    2858 serving.go:312] Generated self-signed cert in-memory
failed to create listener: failed to listen on 0.0.0.0:10251: listen tcp 0.0.0.0:10251: bind: address already in use
Process 2858 has exited with status 1
Detaching
dlv dap (2749) exited with code: 0
```

进源码调试了下, 没仔细看, 还是在网上搜了一下.

## 解决方法

按照参考文章1所说, 在指定了`--config`参数后, 将忽略如下参数.

```log
    --config string            The path to the configuration file. The following flags can overwrite fields in this file:
    --address string   DEPRECATED: the IP address on which to listen for the --port port (set to 0.0.0.0 or :: for listening in all interfaces and IP families). See --bind-address instead. This parameter is ignored if a config file is specified in --config. (default "0.0.0.0")
    --port int         DEPRECATED: the port on which to serve HTTP insecurely without authentication and authorization. If 0, don't serve plain HTTP at all. See --secure-port instead. This parameter is ignored if a config file is specified in --config. (default 10251)
```

...但是人家说的是 1.21 版本, 我的版本是 1.17.2, 启动参数里可没这么说.

```
Insecure serving flags:
    --address string
        DEPRECATED: the IP address on which to listen for the --port port (set to 0.0.0.0 for all IPv4 interfaces and :: for all IPv6 interfaces). See
        --bind-address instead. (default "0.0.0.0")
    --port int
        DEPRECATED: the port on which to serve HTTP insecurely without authentication and authorization. If 0, don't serve plain HTTP at all. See
        --secure-port instead. (default 10251)
```

呵呵.

不过, 在`--config`配置文件中指定其他端口后, 的确启动成功了.

```yaml
apiVersion: kubescheduler.config.k8s.io/v1alpha1
kind: KubeSchedulerConfiguration
## ...省略
healthzBindAddress: 0.0.0.0:11251
metricsBindAddress: 0.0.0.0:11251
```
