# kuber-controller revision.2.kcr中存储的内容不完整[kcr]

kube: 1.16.2

前面说到, `StatefulSet`和`DaemonSet`两种资源在被用户修改后, 会自动创建`ControllerRevision`对象作为备份, 可以用于回滚.

但是最近在写一些`operator`的代码的时候发现, 其中存储的内容不太全啊...

以一个 redis 的 sts 为例, yaml 的内容大致如下

```yaml
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: redis-app
spec:
  ## headless service名称
  serviceName: "redis-service"
  selector:
    matchLabels:
      app: redis
      appCluster: redis-cluster
  replicas: 3
  template:
    metadata:
      labels:
        app: redis
        appCluster: redis-cluster
    spec:
      containers:
      - name: redis
        image: redis
        imagePullPolicy: IfNotPresent
        env:
        - name: name
          value: redis
```

但是生成的`sts`对象和`kcr`的对象如下

![](https://gitee.com/generals-space/gitimg/raw/master/7f76835193f22b595beb87ab94912aa9.png)

左侧`kcr`的`data.spec`下的内容与右侧`sts`的相比, 少了`podManagementPolicy`, `replicas`, `selector`和`serviceName`这几个极为重要的字段.

我想这应该和`StatefulSet`, `DaemonSet`的回滚机制有关, 回滚时不能修改这些字段.

但是我们在使用`kcr`进行备份与回滚的时候, 需要注意.
