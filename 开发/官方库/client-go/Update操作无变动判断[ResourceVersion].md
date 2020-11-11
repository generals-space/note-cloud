# Update操作无变动判断[ResourceVersion]

我们在对资源进行`Update`操作时, 很多时候都要判断其是否真的发生了变化, 如果没有变动, 就不必进行后续的一系列操作.

比如, 如果更新了`Deployment`资源, 但`ResourceVersion`没有发生变化, 说明更新的内容与原本的内容是完全一致的, 那么之后可能就不需要删除重建该`Deployment`名下的`Pod`资源了.

判断变动的场景大概有两处

## 1. AddEventHandler UpdateFunc() 处理函数

```go
	podGroupInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: controller.enqueuePodGroup,
		UpdateFunc: func(old, new interface{}) {
			oldPodGroup := old.(*testgroupv1.PodGroup)
			newPodGroup := new.(*testgroupv1.PodGroup)
			if oldPodGroup.ResourceVersion == newPodGroup.ResourceVersion {
				//版本一致，就表示没有实际更新的操作，立即返回
				return
			}
			controller.enqueuePodGroup(new)
		},
		DeleteFunc: controller.enqueuePodGroupForDelete,
	})
```

## 2. 调用client-go库中的`Update()`方法后

```go
func (c *Controller) processNextWorkItem() bool {
    // ...省略
	sts, err := c.kubeclientset.AppsV1().StatefulSets("kube-system").Get(
		"redis-app",
		apimMetav1.GetOptions{},
	)
	if err != nil {
		klog.Error("================= get sts error %s", err)
		return true
	}
	klog.Infof("sts resource version: %s", sts.ResourceVersion)
	newSTS, err := c.kubeclientset.AppsV1().StatefulSets("kube-system").Update(sts)
	if err != nil {
		klog.Error("================= update sts error %s", err)
		return true
	}
    klog.Infof("new sts resource version: %s", newSTS.ResourceVersion)
    // 此处可以判断双方的`ResourceVersion`值.
	return true
}
```
