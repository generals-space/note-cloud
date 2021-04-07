# builder.update事件处理

参考文章

1. [Support for old object access](https://github.com/kubernetes-sigs/kubebuilder/issues/37)
    - 关于更新事件中, 如何获取旧版本对象信息的讨论
2. [Support for old object access](https://github.com/kubernetes-sigs/kubebuilder/issues/877)

在使用 code-generator 生成的工程中, 我们知道有3种事件处理函数, 如下

```go
	podGroupInformer.Informer().AddEventHandler(cgCache.ResourceEventHandlerFuncs{
		AddFunc: controller.enqueuePodGroup,
		UpdateFunc: func(old, new interface{}) {
			oldPodGroup := old.(*crdV1.PodGroup)
			newPodGroup := new.(*crdV1.PodGroup)
			if oldPodGroup.ResourceVersion == newPodGroup.ResourceVersion {
				//版本一致, 就表示没有实际更新的操作, 立即返回
				return
			}
			controller.enqueuePodGroup(new)
		},
		DeleteFunc: controller.enqueuePodGroupForDelete,
	})

```

但是 kubebuilder 生成的工程, `Reconcile()`函数中却没有任何提示.

```go
func (r *MyClusterReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	_ = context.Background()
	_ = r.Log.WithValues("mycluster", req.NamespacedName)

	// your logic here

	return ctrl.Result{}, nil
}
```

创建/删除这两种操作都还好说, 但是却没有像`UpdateFunc()`那样, 能同时传入`old`, `new`两个版本的对象的方法.

我尝试搜索了一下, 发现这是官方故意设计成这样的. 开发者不需要知道旧版本的对象内容, 如果真的出现必须得到旧对象的时候, 说明代码编写得有问题, 不是`self-healing(自愈)`的系统...

我想了想, 在我目前写过的 operator 代码中, 其实也没有需要过旧对象...那先这样吧...
