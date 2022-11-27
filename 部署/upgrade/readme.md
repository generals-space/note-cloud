# 升级流程

参考文章

1. [kubernetes集群版本升级攻略](https://blog.51cto.com/newfly/2440901)
    - 超全超详细
2. [kubernetes calico IPV6支持](https://www.jianshu.com/p/e92dec9f9cf4)
    - 二进制 + systemd服务脚本升级

## 

1. 先确定要升级的目标版本, 如`v1.17.2`, 然后可以编写相应的`kubeadm-config.v1.17.2.yaml`文件;
2. 将 control plane 节点的`kubeadm`升级到目标版本, 同时可能需要升级同节点上的`kubelet`, `kubectl`;
    - yum/apt这些可能会在升级`kubeadm`的时候连带着将`kubelet`升级到最新版, 所以最好同时指定版本一起升级
3. 使用新版本的`kubeadm`执行`upgrade plan`, 判断目标版本是否合法.
4. `kubeadm upgrade apply --config ./kubeadm-config.v1.17.2.yaml`只会升级 control plane 节点的组件;
5. 升级其余master节点的`kubeadm`(及`kubelet`, `kubectl`), 完成后可能需要重启`kubelet`服务;
6. 升级其余master节点的 kuber 组件, 执行`kubeadm upgrade node experimental-control-plane`
7. 接下来升级 worker 节点, 一台一台来.
8.  先将目标 worker 上的 Pod 驱逐 `kubectl drain $NODE --ignore-daemonsets`
9.  升级 worker 的 `kubeadm`(和`kubelet`, `kubectl`);
10. 拷贝上面编写的`kubeadm-config.v1.17.2.yaml`配置到目标 worker 节点;
11. 在目标 worker 节点上执行`kubeadm upgrade node --config ./kubeadm.v1.17.2.yaml`.
12. 升级成功后需要重启一下 `kubelet`, 然后将该 worker 恢复调度 `kubectl uncordon $NODE`

```
## 安装指定版本的组件
yum install -y kubelet-1.17.2 kubeadm-1.17.2 kubectl-1.17.2
```

