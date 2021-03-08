# kuber-Pod状态[running pending]

`Pod.Status.Conditions`是一个数组, 但貌似没有顺序, 可以简单通过`Condition`成员的`lastTransitionTime`属性自行排序.

PodScheduled ->

## Pending - PodScheduled

`kubectl`查看 Pod 在 Pending 状态时, Pod 的 status 状态如下, `PodScheduled`表示正在被调度(选择主机), `phase`为`Pending`

```yaml
status:
  conditions:
  - lastProbeTime: null
    lastTransitionTime: "xxxxxxx"
    status: "True"
    type: PodScheduled
  phase: Pending
```

`kubectl`查看 Pod 在 Init 状态时, Pod 的 status 状态存在两部分, 一部分是`initContainer`的, 一部分是`container`的.

`kubectl`的`Init:PodInitializing`就表示 initContainer 在执行(如果没有 initContainer, 这个阶段可能会一闪而没).

当然, 当 initContainer 出错时, `Init`后就会出现如`Error`, `CrashLoopBackOff`等状态. 此时 container 的状态可能会出现

```yaml
status:
  containerStatuses:
  - image: xxx
    state:
      waiting:
        reason: PodInitializing
```

> `containerStatuses`和`conditions`平级, 都会存在.

## Running

当 Init 完成(如果有的话), Pod status 中的`initContainerStatuses`都会是`Completed`.

另外不管 container 是`Running`, `Error`还是`CrashLoopBackOff`, Pod status 中的 phase 都还是`Running`. 不过这些状态都会出现在`status.containerStatuses`列表下.

## Terminating - xxx

当通过`kubectl delete`删除 Pod 时, 会出现`Terminating`状态, k8s 会将此 Pod 的地址从 Endpoints 中移除, 不再让新请求进入. 

但是并不会影响 Pod status 中的`phase`状态, 不管 Pod 是在 Running, Pending, 在 kubectl 出现`Terminating`时, 并不会发生改变.

直到 Pod 被真正删除了, 就什么都查不到了.

----

所以