# calico组件启动失败-failed to setup network for sandbox “xxx“：plugin type=“calico“ failed (add)：error getting ClusterInformation：Get “...default“：EOF

参考文章

1. [准备 Kubernetes 集群环境](https://george.betterde.com/cloud-native/20220928.html)

- kube: v1.22.0
- calico: v1.22.4

```log
Events:
  Type     Reason                  Age                 From               Message
  ----     ------                  ----                ----               -------
  Normal   Scheduled               18m                 default-scheduler  Successfully assigned calico-system/calico-kube-controllers-7b7d58c6fc-w7k8v to k8smaster
  Warning  FailedCreatePodSandBox  16m                 kubelet            Failed to create pod sandbox: rpc error: code = Unknown desc = failed to setup network for sandbox "3b547fef6481b6f4f72ffa5d6c9f1be0041e61d83d10ca60a63955084d35b849": plugin type="calico" failed (add): error getting ClusterInformation: Get "https://10.96.0.1:443/apis/crd.projectcalico.org/v1/clusterinformations/default": EOF
```

## 解决方法

因为给 Containerd 配置了代理，导致启动的容器也无法正常访问 Kubernetes 的 Service IP，禁用 Containerd 代理，再重启 Containerd 就正常了。

