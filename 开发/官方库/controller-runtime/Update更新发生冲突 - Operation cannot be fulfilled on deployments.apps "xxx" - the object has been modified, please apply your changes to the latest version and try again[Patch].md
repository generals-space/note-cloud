# Update更新发生冲突 - Operation cannot be fulfilled on deployments.apps "xxx" - the object has been modified, please apply your changes to the latest version and try again[Patch]

参考文章

1. [Object Has Been Modified Error After Upgrading to v0.9.0-beta.0](https://github.com/kubernetes-sigs/controller-runtime/issues/1509)
    - This is a thing that happens if someone else updates the object **between your get and update**. 
    - Kubernetes uses optimistic locking, and the **Get is on an async cache**
    - The way to fix these more generally is to switch from Update() to Patch() calls, preferably server-side-apply if possible but failing that a MergeFrom patch might suffice.
2. [How to elegantly solve the update conflict problem](https://github.com/kubernetes-sigs/controller-runtime/issues/1748)
    - Better solution, stop using Update(). I can't think of any reason to use it in controller code. 
    - Most requests should use Patch() with Server Side Apply and the few that can't use Apply should use auto-merge-patch.

## 问题描述

在controller中调用`r.Update()`更新`Deployment`类型资源时, 频繁出现下面的错误.

```
2022-03-12T18:59:33.966+0800	ERROR	controller-runtime.controller	Reconciler error	{"controller": "deployment", "request": "test-test/0311-07", "error": "Operation cannot be fulfilled on deployments.apps \"0311-07\": the object has been modified; please apply your changes to the latest version and try again"}
```

## 解决方案

其实在编写 crd controller 时, 偶尔也遇到过这个问题, 但是没有像这次一样这么频繁. 

按照参考文章1中的说法, 这是因为在我们代码逻辑中, `Get()`与`Update()`之间, 此资源对象被其他地方修改过了, `resourceVersion`字段发生了改变, 所以更新才会失败.

我们自己编写的 crd controller 是代码中自己控制的, 但是对 deployment 等原生资源的改动, 却是会受到原生 controller 的影响的, 所以发生的频次比 crd controller 中又多很多.

------

当然, 可以`Watch`到这样的变动, 更新失败返回错误后, 此事件变动还是会再次进入`Reconcile()`, 总有一次可以成功的.

但是这样做其实是非常不"优雅"的, 建议 controller 使用`Patch()`代替`Update()`, 