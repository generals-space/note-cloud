参考文章

1. [Installing calicoctl as a Kubernetes pod](https://docs.projectcalico.org/v3.10/getting-started/calicoctl/install#installing-calicoctl-as-a-kubernetes-pod)
    - calicoctl 安装手册
2. [calicoctl user reference](https://docs.projectcalico.org/v3.10/reference/calicoctl/)
    - calicoctl 使用手册
3. [Calico网络方案](https://www.cnblogs.com/netonline/p/9720279.html)
    - calico二进制文件部署方案
    - calicoctl使用方法
    - flannel与calico网络方案的简单对比

## 1. 安装

官方提供了3种calicoctl的安装方式: 

1. 二进制文件(裸机部署)
2. docker镜像
3. 作为kuber中的一个pod.

其实这三种方式都算是同一种, 二进制文件需要读取配置文件, 运行的pod也是需要sa等相关权限的支持, 主要目的就是要从etcd或是通过kuber API获取网络状态的信息.

## 使用

`calicoctl get node`: 查看集群中的节点信息

`calicoctl get ippool`: 查看集群中pod网段范围(这个值应该是在`calico-node`的daemonset的部署文件中, 由`CALICO_IPV4POOL_CIDR`字段定义的, 而且应该是与集群的apiserver对pod网段配置是相同的). 与`kubectl get ippool`结果相同, 通过calicoctl对ippool的CURD操作应该与kubectl的功能是相同的.

`calicoctl node status`: 显示邻居节点的信息

