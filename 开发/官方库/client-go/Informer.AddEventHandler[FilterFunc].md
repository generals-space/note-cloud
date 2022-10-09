# 

```go
import (
	"k8s.io/client-go/tools/cache"
)

setInformer.Informer().AddEventHandler(
		cache.ResourceEventHandlerFuncs{
			DeleteFunc: ssc.deleteStatefulSet,
		},
	)
```

```go
import (
	"k8s.io/client-go/tools/cache"
)
	podInformer.Informer().AddEventHandler(
		cache.FilteringResourceEventHandler{
			// 只有 FilterFunc 返回 true 的 Pod 才会进入 Handler 阶段
			FilterFunc: func(obj interface{}) bool {
				switch t := obj.(type) {
				case *v1.Pod:
					return assignedPod(t)
				case cache.DeletedFinalStateUnknown:
					if pod, ok := t.Obj.(*v1.Pod); ok {
						return assignedPod(pod)
					}
					return false
				default:
					return false
				}
			},
			Handler: cache.ResourceEventHandlerFuncs{
				AddFunc:    
				UpdateFunc: 
				DeleteFunc: 
			},
		},
	)
```
