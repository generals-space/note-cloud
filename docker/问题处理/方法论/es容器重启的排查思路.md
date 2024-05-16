容器重启, 而非pod被删, 引发了es集群中的节点脱离集群.

使用 kubectl get pod 可以看到容器的资源大于0, 进而使用 kubectl logs -p pod名称, 查看容器重启前发生了什么, 尤其是得到最后一条日志发生的时间.

```log
Caused by: org.elasticsearch.common.util.concurrent.EsRejectedExecutionException: rejected execution of org.elasticsearch.transport.TransportService$7@6f3e918 on EsThreadPoolExecutor[search, queue capacity = 5000, org.elasticsearch.common.util.concurrent.EsThreadPoolExecutor@430f230[Running, pool size = 13, active threads = 13, queued tasks = 5000, completed tasks = 8773540]]
```

一般来说, es这种容器异常重启的很大可能是OOM了, 登陆该容器所在主机, 查看/var/log/message文件, 在最后一条日志的时间范围内, 是否存在oom日志, 结果没有.

然后使用 kubectl describe pod 查看pod的事件, 是否是因为CPU占用过高, liveness检查没通过, 被 kubelet 干掉了. 但是中间件没有配置 liveness 探针, 就算健康检查没通过, 容器也只会保留在原有的状态, 不会被干掉.

进而回到日志本身, 除去上面2种可能后, 那么就只剩下ES进程自身退出引发的容器重启了.
