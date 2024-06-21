# kubelet 无法正常运行 - Orphaned pod found, but volume paths are still present on disk

参考文章

1. [Orphaned pod found - but volume paths are still present on disk](https://github.com/kubernetes/kubernetes/issues/60987)
2. [Orphaned pod found, but volume paths are still present on disk](https://github.com/mattshma/bigdata/issues/114)
    - 提示风险, 值得注意

centos: 7
kube: v1.13.2

kuber 集群中某个宿主机异常宕机后, kubelet 无法正常运行, docker ps 没也没容器启动. 查看`/var/log/message`, 发现不停有如下错误日志.

```log
kubelet: E0309 16:46:30.429770 3112 kubelet_volumes.go:128] Orphaned pod "2815f27a-219b-11e8-8a2a-ec0d9a3a445a" found, but volume paths are still present on disk : There were a total of 1 errors similar to this. Turn up verbosity to see them.
```

按照参考文章1中所说, 到`/var/lib/kubelet/pods`下把日志中的孤儿 pod 目录删除即可, 我自己删除后可以解决这个问题(虽然后面还出现了其他问题, 但已经和这篇文章的主题关系不大了).

不过参考文章2中提到了删除孤儿 pod 目录的风险, 目前还没有遇到, 但是值得警惕.
