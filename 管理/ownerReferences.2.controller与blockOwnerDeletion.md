参考文章

1. [为什么kubernetes的ownerReference要设计成列表，存在多个ownerReference的情况吗？](https://joshua.su/8a98cc78605e)
2. [kubernetes/design-proposals-archive](https://github.com/kubernetes/design-proposals-archive/blob/acc25e14ca83dfda4f66d8cb1f1b491f26e78ffe/api-machinery/controller-ref.md)
    - 官方文档
3. [k8s中的OwnerReference和Finalizers](https://blog.csdn.net/qq_50519011/article/details/136184126)
    - k8s中的级联删除策略
4. [k8s ownerReferences](https://zhuanlan.zhihu.com/p/577185348)
    - orphan解除级联关系: 删除属主对象后，会忽略附属资源，附属资源会删除`ownerReferences`

前面提到过, `ownerReference`的级联删除只需要4个字段就可以实现.

- apiVersion
- kind
- name
- uid

还剩下`blockOwnerDeletion`与`controller`的作用未明.

## blockOwnerDeletion

要认识这个字段, 需要先了解k8s中的级联删除策略.

**k8s中的级联删除策略**

- background: 先删除属主资源, 再在后台删除从属资源(默认策略)
    - 使用 kubectl delete 命令, 以及使用 client-go 提供的`Delete()`函数删除资源, 默认使用`background`.
- foreground: 先删除从属资源, 再删除属主资源(一般与BlockOwnerDeletion=true结合使用)
- orphan: 不考虑OwnerReference, 只删除该资源, 不级联删除
    - 很容易理解, 需要注意的是, 当父资源使用`orphan`策略删除时, 子资源中相应的`ownerReference`记录会被自动移除.

`blockOwnerDeletion`只在`foreground`删除中会起作用, 不影响`background`与`orphan`策略.

采用`foreground`删除策略时, 如果`blockOwnerDeletion==false(默认)`, 则删除主资源时, 虽然说是先删除从资源, 但并不会等待从资源真的被删除(即异步删除), 而是直接返回删除结果.

而如果`blockOwnerDeletion==true`, 则删除主资源时, 会先等待从资源真的被删除后, 再返回. 在使用 client-go Delete() 函数时, 会阻塞直到子资源与父资源被删掉之后才会返回.

按照参考文章4中所说, `foreground`删除主资源时, 会有如下变化

- 对象仍然可以通过 REST API 可见
- 会设置对象的`deletionTimestamp`字段
- 对象的`finalizers`字段包含了值`foregroundDeletion`

如果没有从资源拥有`blockOwnerDeletion==true`配置, 那么上述状态可能会一闪而逝, 毕竟主资源马上就会被删掉了.

## controller

这个字段的设计目的, 是为了满足`ownerReferences[]`的列表式实现.

我们知道, deployment资源是通过`ReplicaSet`管理`Pod`的, `Pod`的owner就是`ReplicaSet`. 

但`ownerReferences`是一个列表, 如果开发者在自己编写的Operator中, 将自已定义的CRD资源也写到这种Pod的`ownerReferences`中, 会发生什么?

控制权此时发生了撞车, `ReplicaSet`会不会撒手不管了?

这就是`ownerReference.controller`的作用, 按照参考文章2中的官方提议, 该字段规定了, 在所有的`ownerReference`中, 只有第一个进行申请的owner才能将`controller`字段置为true. 当其它服务发现当前资源已经有了一个存在的`ContorllerRef`时其无法再将自己的`ControllerRef`置为`true`, 并且也**不应该**再把当前资源算作自己期望状态的一部分.

------

关于这一点, 在许多官方库中也有体现. 

比如, [[controller-runtime]:pkg/handler/enqueue_owner.go](https://gitee.com/skeyes/controller-runtime/blob/438d738ad99b7c08e1a3fd22fe37764cb7bf7e61/pkg/handler/enqueue_owner.go#L159)

开发者在关注自定义 CRD 资源的同时, 也可以通过`Owns()`函数关注 CRD 的子资源.

```go
func (r *RedisClusterReconciler) SetupWithManager(ctx context.Context, mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1.RedisCluster{}).
		Owns(&corev1.Pod{}).
		Complete(r)
}
```

`EnqueueRequestForOwner.getOwnerReconcileRequest()`在监听子资源的时候, 会遍历ta的`ownerReferences`数组, 判断自己的 CRD 类型是否抢到了子资源的 controller 标记.

但是, `EnqueueRequestForOwner`也提供了选项, 就算是没有抢到 controller 标记, 也可以去做一些操作.

看起来, 全靠自觉啊...

同样, [[apimachinery]:pkg/apis/meta/v1/controller_ref.go](https://e.gitee.com/skeyes/repos/skeyes/apimachinery/blob/52c7025bffab65d98068d11d6c222ce5dc42a2b3/pkg/apis/meta/v1/controller_ref.go#L33) 提供的`GetControllerOf()`工具函数, 获取的目标子资源真正的属主.
