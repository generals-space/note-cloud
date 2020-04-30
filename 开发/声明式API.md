# 声明式API

参考文章

1. [【Kubernetes】声明式API与Kubernetes编程范式](https://www.cnblogs.com/yuxiaoba/p/9801161.html)
    - `kubectl apply` vs `kubectl replace`, 前者即为声明式请求, 因为ta可以同时处理多个任务, 且具备`merge`能力.
    - 以`istio`将`envoy`容器注入到Pod的执行过程为例, 介绍了在 kuber 集群中通过声明式API更新对象的能力(不过感觉不太适合这个场景, 倒是`istio`的工作流程讲解得很清楚)
    - **声明式API才算Kubernetes项目编排能力赖以生存的核心所在**
2. [【Kubernetes】深入解析声明式API](https://www.cnblogs.com/yuxiaoba/p/9803284.html)
    - client提交创建资源的请求后, apiserver的执行流程(这个过程写的很详细, 就是感觉有点不符合"声明式API"的主题)
