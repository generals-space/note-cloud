# kubectl查看initContainers日志

参考文章

1. [调试 Init 容器](https://kubernetes.io/zh/docs/tasks/debug-application-cluster/debug-init-containers/)

如果一个 yaml 部署文件中包含`initContainers`, 在pod启动后使用`k logs -f`只能看到`containers`容器的日志, 看不到初始化容器的(而且初始化阶段没法看日志).

想要查看`initContainers`中容器的日志, 可以使用`-c`参数.

```
kubectl logs <pod-name> -c <init-container-2>
```
