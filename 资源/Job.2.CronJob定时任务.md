# kuber-Job.2.CronJob定时任务

参考文章

1. [Kubernetes核心概念总结](https://www.cnblogs.com/zhenyuyaodidiao/p/6500720.html)
2. [官方文档 - Running automated tasks with cron jobs](https://kubernetes.io/docs/tasks/job/automated-tasks-with-cron-jobs/)
3. [TTL Mechanism for Finished Jobs](https://kubernetes.io/docs/concepts/workloads/controllers/jobs-run-to-completion/#ttl-mechanism-for-finished-jobs)
    - `ttlSecondsAfterFinished`设置job结束后多久会自动删除.

`CronJob`的`schedule`可以是linux常用的cron任务形式, `container`也可以是普通容器. 只不过执行完成就被销毁, 等待下次执行.

```yaml
apiVersion: batch/v1beta1
kind: CronJob
metadata:
  name: hello
spec:
  schedule: "*/1 * * * *"
  jobTemplate:
    spec:
      ttlSecondsAfterFinished: 60
      template:
        spec:
          containers:
          - name: hello
            image: busybox
            args:
            - /bin/sh
            - -c
            - date; echo Hello from the Kubernetes cluster
          restartPolicy: OnFailure
```

创建的cronjob并不会立刻创建pod容器, 而是在第一次触发的时间开始创建并执行第一次任务.

```log
$ k get cronjob
NAME    SCHEDULE      SUSPEND   ACTIVE   LAST SCHEDULE   AGE
hello   */1 * * * *   False     0        <none>          16s
```

之后该 cronjob 会自动创建 Job 资源对象.

```log
$ k get job
NAME               COMPLETIONS   DURATION   AGE
hello-1594731900   1/1           18s        98s
hello-1594731960   1/1           24s        38s
```

一个`*/30 * * * *`的job, 状态变化如下

```log
$ k get pod
hello-1594731900-vsgj9                  0/1     Completed   0          2m18s
hello-1594731960-8kz78                  0/1     Completed   0          78s
hello-1594732020-5cz89                  0/1     Completed   0          18s
```
可以看到, 30秒后pod状态变为`Complete`, 并不会删除(`ttlSecondsAfterFinished`没有生效.).

可以使用`kubectl get job`获取执行过的任务信息.
