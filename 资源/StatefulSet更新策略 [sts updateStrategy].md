# sts-更新策略

参考文章

1. [Kubernetes资源对象：StatefulSet](https://blog.csdn.net/fly910905/article/details/102092570)
    - 应用场景: 稳定的持久化存储, 稳定的网络标识, 有序部署与有序收缩.
    - 更新策略, 解释了`.spec.updateStrategy.rollingUpdate.partition`的作用
2. [Kubernetes指南 StatefulSet](https://feisky.gitbooks.io/kubernetes/concepts/statefulset.html)
  - 给出了更新策略中的`partition`和管理策略中的`parallel`的使用示例.

`statefulset`目前支持两种策略

- `OnDelete`: 当`.spec.template`更新时, 并不立即删除旧的`Pod`, 而是等待用户手动删除这些旧`Pod`后自动创建新Pod. 这是默认的更新策略, 兼容 v1.6 版本的行为.
- `RollingUpdate`: 当`.spec.template`更新时, 自动删除旧的`Pod`并创建新`Pod`替换. 在更新时, 这些`Pod`是按逆序的方式进行, 依次删除、创建并等待`Pod`变成`Ready`状态才进行下一个`Pod`的更新. 

其中`RollingUpdate`有一个`partition`选项, 只有序号大于或等于`partition`的`Pod`会在`.spec.template`更新的时候滚动更新, 而其余的`Pod`则保持不变(即便是删除后也是用以前的版本重新创建).

```yaml
spec:
  ## headless service名称
  serviceName: "redis-service"
  replicas: 6
  updateStrategy:
    type: RollingUpdate
    rollingUpdate:
      partition: 4
```

> `partition`从0开始计数

这样, 在更新images版本后apply, 你会发现只有`redis-app-4`和`redis-app-5`会更新, 其他的`Pod`则保持不动.
