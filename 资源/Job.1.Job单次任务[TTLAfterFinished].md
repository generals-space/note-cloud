# kuber-Job.1.Job单次任务

参考文章

1. [Kubernetes Job配置](https://www.cnblogs.com/breezey/p/6582754.html)
2. [TTL Mechanism for Finished Jobs](https://kubernetes.io/docs/concepts/workloads/controllers/jobs-run-to-completion/#ttl-mechanism-for-finished-jobs)
    - `ttlSecondsAfterFinished`设置job结束后多久会自动删除.
3. [k8s 自动清理已完成的Job相关的资源对象（v1.12版本以上才支持）](https://blog.csdn.net/H12590400327/article/details/88647429)
    - `ttlSecondsAfterFinished`特性需要`apiserver`, `controller`和`scheduler`开启相应的特性: `TTLAfterFinished`
    - 3个组件的yaml配置中都添加`--feature-gates=TTLAfterFinished=true`

`Job`资源用于处理一些一次性任务, 比如初始化工作(好像可以用`initContainers`), 或是删除后的数据清理工作(貌似也可以用`preStop`完成)...呃...

总感觉不太实用的样子, 很容易使用其他方法代替.

在Job的定义中, `restartPolicy`(重启策略)只能是`Never`和`OnFailure`. 

Job可以控制一次性任务的Pod的完成次数(`Job.spec.completions`)和并发执行数(`Job.spec.parallelism`), 当Pod成功执行指定次数后, 即认为Job执行完毕. 

```yaml
apiVersion: batch/v1
kind: Job
metadata:
  labels:
    name: myjob
  name: myjob
spec:
  ttlSecondsAfterFinished: 60
  template:
    metadata:
      name: myjob
    spec:
      containers:
      - name: centos7
        image: registry.cn-hangzhou.aliyuncs.com/generals-space/centos:7
        command:
        - sleep
        - "300"
      restartPolicy: Never
```

在 command 执行过程中, job 状态如下.

```log
$ k get job
NAME    COMPLETIONS   DURATION   AGE
myjob   0/1           15m        15m
```

完成后则如下

```log
$ k get job
NAME    COMPLETIONS   DURATION   AGE
myjob   1/1           15m        16m
```

`300`秒后该`Job`就会自动删除.

参考文章3说这个功能需要开启3个组件的`TTLAfterFinished`特性, 不过这个特性从`1.12`就是`ALPHA`, 我现在测试的集群版本是`1.16`了, 还是`ALPHA`...可真牛p.

在未开启该特性时, 就算在yaml中指定了`ttlSecondsAfterFinished`字段, 最终生成的job资源中也不会出这个字段出现.
