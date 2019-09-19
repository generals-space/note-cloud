# kuber-Job之定时任务

参考文章

1. [Kubernetes核心概念总结](https://www.cnblogs.com/zhenyuyaodidiao/p/6500720.html)

2. [官方文档 - Running automated tasks with cron jobs](https://kubernetes.io/docs/tasks/job/automated-tasks-with-cron-jobs/)

3. [TTL Mechanism for Finished Jobs](https://kubernetes.io/docs/concepts/workloads/controllers/jobs-run-to-completion/#ttl-mechanism-for-finished-jobs)
    - `ttlSecondsAfterFinished`设置job结束后多久会自动删除.


从程序的运行形态上来区分，我们可以将Pod分为两类：长时运行服务(jboss、mysql等)和一次性任务(数据计算、测试). RC创建的Pod都是长时运行的服务，而Job创建的Pod都是一次性任务. 

在Job的定义中，`restartPolicy`(重启策略)只能是`Never`和`OnFailure`. Job可以控制一次性任务的Pod的完成次数(`Job.spec.completions`)和并发执行数(`Job.spec.parallelism`)，当Pod成功执行指定次数后，即认为Job执行完毕. 

网上有很多介绍关于`CronJob`的文章, 但大多数还停留在`batch/v2alpha1`的api版本上, 创建cronjob时会报如下错误.

```
$ kubectl create -f cronjob.yaml 
error: unable to recognize "cronjob.yaml": no matches for kind "CronJob" in version "batch/v2alpha1"
```

所以还是按照官方的文档更准确.

`CronJob`的`schedule`可以是linux常用的cron任务形式, `container`也可以是普通容器. 只不过执行完成就被销毁, 等待下次执行.

```yml
apiVersion: batch/v1beta1
kind: CronJob
metadata:
  name: hello
spec:
  schedule: "*/1 * * * *"
  jobTemplate:
    spec:
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

一个`*/30 * * * *`的job, 状态变化如下

```
wrk-job-1567591800-dc59z                  1/1     Running             0          35s    192.168.171.25    k8s-worker-7-18   <none>           <none>
wrk-job-1567591800-dc59z                  0/1     Completed           0          65s    192.168.171.25    k8s-worker-7-18   <none>           <none>
```
可以看到, 30秒后pod状态变为Complete, 并不会删除.

可以使用`kubectl get job`获取执行过的任务信息.
