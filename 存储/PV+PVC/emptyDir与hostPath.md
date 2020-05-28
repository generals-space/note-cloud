# emptyDir与hostPath

参考文章

1. [emptyDir与hostPath](https://www.cnblogs.com/breezey/p/9827570.html)

`emptyDir`是k8s从宿主机上随机选取一个目录, 无需手动指定路径, 且当Pod被移除时, 该目录中的数据将被移除.

`hostPath`需要用户手动指定目标路径, Pod移除后不会被删除.
