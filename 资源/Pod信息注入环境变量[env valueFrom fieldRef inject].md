# kuber-把pod本身的信息注入到pod中

参考文章

1. [Expose Pod Information to Containers Through Environment Variables](https://kubernetes.io/docs/tasks/inject-data-application/environment-variable-expose-pod-information/)
    - 通过env的形式挂载pod/container信息字段
2. [Expose Pod Information to Containers Through Files](https://kubernetes.io/docs/tasks/inject-data-application/downward-api-volume-expose-pod-information/)
    - 通过volume的形式挂载这些字段
3. [Capabilities of the Downward API](https://kubernetes.io/docs/tasks/inject-data-application/downward-api-volume-expose-pod-information/#capabilities-of-the-downward-api)
    - 貌似这里是所有可用的字段

有时程序需要知道自己运行的pod名称, 所在的namespace, 或是运行在哪个node节点上, 这些对于`kube-system`这种管理型的程序比较有用(比如flannel).

但这些信息在Pod被调度然后成功启动后才能确定下来, 如果有办法

下面节选自flannel的部署文件

```yaml
apiVersion: apps/v1
kind: DaemonSet
metadata:
  name: kube-flannel-ds-amd64
  namespace: kube-system
  labels:
    app: flannel
spec:
    ## ...省略
    spec:
      containers:
      - name: kube-flannel
        image: quay.io/coreos/flannel:v0.11.0-amd64
        env:
        - name: POD_NAME
          valueFrom:
            fieldRef:
              fieldPath: metadata.name
        - name: POD_NAMESPACE
          valueFrom:
            fieldRef:
              fieldPath: metadata.namespace
        - name: NODE_NAME
          valueFrom:
            fieldRef:
              fieldPath: spec.nodeName
        - name: POD_IP
          valueFrom:
            fieldRef:
              fieldPath: status.podIP
```

这样, 在flannel的pod启动后, 可以从环境变量里得到`POD_NAME`, `POD_NAMESAPCE`, `NODE_NAME`等信息.

其实可注入的就是`metadata`, `spec`, `status`这几部分的内容, 使用`k get pod pod名称 -o yaml`, 这几块下面的字段都可以尝试下. 

比如想获得Pod所在宿主机的IP地址, 可以使用`status.hostIP`来确定.
