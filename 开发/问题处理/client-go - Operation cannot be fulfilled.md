# client-go - Operation cannot be fulfilled

参考文章

1. [please apply your changes to the latest version and try again](https://github.com/kubernetes/kubernetes/issues/84430)

k8s 版本: 1.13.2

场景描述

我自己编写的 crd, 在进行 update 操作的时候, 更新失败, 报错显示如下

```
Operation cannot be fulfilled on nodes "192.168.0.1": the object has been modified; please apply your changes to the latest version and try again
```

> `on nodes "192.168.0.1"`是我瞎写的, 实际是待更新的目标 crd 资源对象.

参考文章1中有人回答说是 k8s 的 bug, 吓了我一跳(ta的版本是1.13.6, 非常相近).

但是最终经过我的排查, 是因为我在 update 的处理函数中, 又调用了一次 crd 的`Update()`方法.

```go
	podGroupInformer.Informer().AddEventHandler(cgCache.ResourceEventHandlerFuncs{
		AddFunc: controller.enqueuePodGroup,
		UpdateFunc: func(old, new interface{}) {
            // 一般来说, update 的处理函数都是对 new 对象进行操作的.

            // 对更新过的 new 对象再做一点修改, 一般是 status 部分.
            new.(podGroup).xxx
            // 这里嵌套了...
            controller.Update(new);
		},
		DeleteFunc: controller.enqueuePodGroupForDelete,
	})
```

比如在某次`UpdateFunc`操作中, 我们调用的`Update()`方法将原来版本为1的 crd 资源更新到了版本2. 

更新成功后, informer 接收到了这个 update 事件, 再次进入到`UpdateFunc`, 此时传入的参数中, `old`的版本为1, `new`的版本为2, 仍是对 new 做一些修改, 再次进行`Update`, 但是目标的版本还是`new`的2, 要知道, 此时(第2次进入`UpdateFunc`) k8s 中已经存在一个版本为 2 的 crd 资源了, 所以这次的更新会出错.

有点乱, 总之就是嵌套了...

但是总比无限递归要好.
