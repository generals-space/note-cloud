参考文章

1. [Configure a Security Context for a Pod or Container](https://kubernetes.io/docs/tasks/configure-pod-container/security-context/)
    - 官方文档
2. [为 Pod 或容器配置安全上下文](https://kubernetes.io/zh-cn/docs/tasks/configure-pod-container/security-context/)
    - 官方文档
    - 参考文章1的中文版
    - `readOnlyRootFilesystem: true`, 以只读方式加载容器的根文件系统(默认为`false`)
    - `securityContext`可以是 Pod 级别, 也可以是 container 级别.


```yaml
apiVersion: v1
kind: Pod
metadata:
  name: security-context-demo
spec:
  securityContext:
    runAsUser: 1000
    runAsGroup: 3000
    fsGroup: 2000
  volumes:
  - name: sec-ctx-vol
    emptyDir: {}
  containers:
  - name: sec-ctx-demo
    image: busybox:1.28
    command: [ "sh", "-c", "sleep 1h" ]
    volumeMounts:
    - name: sec-ctx-vol
      mountPath: /data/demo
    securityContext:
      allowPrivilegeEscalation: false
      readOnlyRootFilesystem: true
```

> 注意: 官网的示例中, `allowPrivilegeEscalation`和`readOnlyRootFilesystem`不能配置到`spec.securityContext{}`块下, 只能配置在`containers[].securityContext{}`块, 否则会报错. 
> 
> 前者只接受`runAsUser`等用户组配置, 而后者可以接受所有配置.
