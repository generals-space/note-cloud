# volumeClaimTemplates - 让statefulset创建的pod每个都有单独的存储

参考文章

1. [为什么deployment创建的pod共享一个存储，statefulset创建的pod每个都有单独的存储](https://blog.csdn.net/MssGuo/article/details/127818345)

## 场景描述

希望sts下的每个Pod挂载外部的不同路径(但是挂载到Pod内部时还是相同的), 如果使用常规的`volumes`字段, 所有Pod将会同用同一个实际路径.

这种场景可以使用`volumeClaimTemplates`实现.

```yaml
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: my-statefulset
spec:
  serviceName: "my-service"
  replicas: 2
  selector:
    matchLabels:
      app: my-app
  template:
    metadata:
      labels:
        app: my-app
    spec:
      containers:
      - name: my-container
        image: registry.cn-hangzhou.aliyuncs.com/generals-space/centos:7
        volumeMounts:
        - name: pvc_name
          mountPath: /opt/log
  volumeClaimTemplates:
  - metadata:
      name: pvc_name
    spec:
      accessModes:
        - ReadWriteMany
      storageClassName: nfs
      resources:
        requests:
          storage: 1Gi
```

k8s会自动创建名为`pvc_name`的pvc, 且只有一个, 但是会出现3个pv资源, 每个Pod绑定的实际pv资源是不同的.

由于`volumeClaimTemplates`中定义的是pvc资源的模板, 因此无法直接用 hostPath/emptyDir, 需要`storageClassName`通过 provisioner 间接完成 pv 的创建.

貌似 deployment 不支持这种方式.
