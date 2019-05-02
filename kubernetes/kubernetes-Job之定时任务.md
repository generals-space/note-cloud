# kubernetes-Job之定时任务

参考文章

1. [Kubernetes核心概念总结](https://www.cnblogs.com/zhenyuyaodidiao/p/6500720.html)

2. [官方文档 - Running automated tasks with cron jobs](https://kubernetes.io/docs/tasks/job/automated-tasks-with-cron-jobs/)

从程序的运行形态上来区分，我们可以将Pod分为两类：长时运行服务（jboss、mysql等）和一次性任务（数据计算、测试）。RC创建的Pod都是长时运行的服务，而Job创建的Pod都是一次性任务。

在Job的定义中，`restartPolicy`（重启策略）只能是`Never`和`OnFailure`。Job可以控制一次性任务的Pod的完成次数（`Job.spec.completions`）和并发执行数（`Job.spec.parallelism`），当Pod成功执行指定次数后，即认为Job执行完毕。

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