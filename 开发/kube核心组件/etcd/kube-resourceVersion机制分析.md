# kube-resourceVersion机制分析

参考文章

1. [Etcd 中 Revision, CreateRevision, ModRevision, Version 的含义](https://www.cnblogs.com/FengZeng666/p/16156407.html)
    - etcd 的 mvcc 机制
    - etcd 的 watch 命令及参数使用方法
2. [Kubernetes-resourceVersion机制分析](https://fankangbest.github.io/2018/01/16/Kubernetes-resourceVersion%E6%9C%BA%E5%88%B6%E5%88%86%E6%9E%90/)
    - k8s etcd 客户端如何处理 resourceVersion 信息

## 关于 watch

```
etcdctl watch '/registry/statefulsets/default/test-sts'
```

watch 操作可以监听任意 key, 即使这个 key 不存在也不会报错.

不会在初始时获取全量数据, 只会获取增量数据, 这一点要与 client-go 的 reflector 实现区分开.

## version

```log
# etcdctl get '/registry/statefulsets/default/test-sts' --write-out="fields" | grep -v Value
"ClusterID" : 6138936026546267994
"MemberID" : 4601696108602649857
"Revision" : 298066
"RaftTerm" : 5
"Key" : "/registry/statefulsets/default/test-sts"
"CreateRevision" : 282406
"ModRevision" : 282410
"Version" : 2
"Lease" : 0
"More" : false
"Count" : 1
```

- Revision: 作用域为集群, 逻辑时间戳, 全局单调递增, 任何 key 的增删改都会使其自增
    - k8s集群中, 每个key更新都会引起该字段自增.
- CreateRevision: 作用域为 key, 等于创建这个 key 时集群的 Revision, 直到删除前都保持不变
- ModRevision: 作用域为 key, 等于修改这个 key 时集群的 Revision, 只要这个 key 更新都会自增
    - key 发生变更时的 Revision, 在 k8s 中, 由于 key 太多, 更新频繁, 所以该值应该是会"跳变"的.
    - 在 k8s 中, 资源的`metadata.resourceVersion`就是这个值.
- Version: 作用域为 key, 这个key刚创建时Version为1, 之后每次更新都会自增, 即这个key从创建以来更新的总次数. 
