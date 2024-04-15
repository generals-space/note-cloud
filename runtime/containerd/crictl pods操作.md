# crictl pods操作

crictl 相比于 docker, 多了一个 pod 视角. 因为 containerd 启动的容器是有 namespace 的, 且所有的容器都需要创建一个 pod 资源(与 kube 的 pod 概念是相同的).

```
[root@k8s-master-01 ~]# crictl pods
POD ID              CREATED             STATE               NAME                                        NAMESPACE       ATTEMPT  RUNTIME
fb769a38a438d       About an hour ago   Ready               kube-controller-manager-k8s-master-01       kube-system     4        (default)
e4ff3966719ba       About an hour ago   Ready               etcd-k8s-master-01                          kube-system     4        (default)
82d113cca6b03       12 hours ago        NotReady            cert-manager-cainjector-76bbdd77f7-bbnwv    cert-manager    0        (default)
455f1dcb08660       25 hours ago        NotReady            cert-manager-cdbc489b6-7kglh                cert-manager    3        (default)
```

## 删除

crictl rmp pod名称
