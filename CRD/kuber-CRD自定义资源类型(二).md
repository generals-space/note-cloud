# kuber-CRD自定义资源类型(二)

参考文章

1. [kubernetes/sample-controller](https://github.com/kubernetes/sample-controller)
2. [coreos/prometheus-operator](https://github.com/coreos/prometheus-operator)

CRD的官方示例非常简单, 实际一点的应用要属参考文章2了. 貌似ingress-controller也是CRD的应用.

今天(20191111)在网银互联面试时, 面试官给出了理想中的CRD的应用场景.

以etcd集群为例, etcd是可以以集群的形式部署在kuber集群中的. 只要写好配置文件, 定义好挂载卷等信息就可以.

但是, 考虑一下, 这样的etcd集群如何升级?

诚然, prometheus官方提供过适用于kuber的部署yaml文件, 但是使用CRD可以更精确地控制特定资源的行为.

比如etcd的滚动升级, 比如mysql升级时主节点升级时, 从节点提升为主节点对外服务等.
