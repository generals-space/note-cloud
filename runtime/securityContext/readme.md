参考文章

1. [Configure a Security Context for a Pod or Container](https://kubernetes.io/docs/tasks/configure-pod-container/security-context/)
    - 官方文档
2. [为 Pod 或容器配置安全上下文](https://kubernetes.io/zh-cn/docs/tasks/configure-pod-container/security-context/)
    - 官方文档
    - 参考文章1的中文版
    - `readOnlyRootFilesystem: true`, 以只读方式加载容器的根文件系统(默认为`false`)
    - `securityContext`可以是 Pod 级别, 也可以是 container 级别.
