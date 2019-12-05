[官方文档 CREATING YOUR OWN CHARTS](https://helm.sh/docs/intro/using_helm/#creating-your-own-charts)
    - 创建新的chart工程(`helm create`)及错误检查(`helm lint`), 打包操作(`helm package`)等
[Helm v3.0.0 安装和使用](https://blog.csdn.net/twingao/article/details/103218363)
    - 个人的使用记录, 给出了3个常用chart源.

## 关于多个源的pull/install操作.

`helm search hub mysql`: 在官方 hub 中搜索 package, 这其实是一个集合, 但并不能直接install/pull, 需要从指定源中操作.

```
$ helm search hub redis
URL                                                     CHART VERSION   APP VERSION     DESCRIPTION
https://hub.helm.sh/charts/bitnami/redis                10.1.0          5.0.7           Open source, advanced key-value store. It is of...
https://hub.helm.sh/charts/stable/redis-ha              4.1.1           5.0.5           Highly available Kubernetes implementation of R...
https://hub.helm.sh/charts/stable/prometheus-re...      3.2.0           1.0.4           Prometheus exporter for Redis metrics
```

> `bitnami`和`stable`是不同的源.

`helm search repo redis`: 如果添加了多个源, 这种搜索方式将查询本地所有的源.

`helm search repo stable/redis`: 在本地安装的指定源中搜索指定 package.

`helm fetch stable/wordpress`: 下载stable源中的chart包

## 常用源

1. `stable`:     https://kubernetes-charts.storage.googleapis.com
2. `bitnami`:    https://charts.bitnami.com/bitnami
3. `aliyuncs`:   https://apphub.aliyuncs.com

> `helm repo add stable https://kubernetes-charts.storage.googleapis.com`即可添加新的源.
