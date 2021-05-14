`informer`都是对于某种特定资源的

## 初始化方式

```go
	// Shared指的是多个 lister 共享同一个cache, 而且资源的变化会同时通知到cache和listers.
	factory := informers.NewSharedInformerFactory(clientset, 0)
	// nodeInformer 拥有两个方法: Informer, Lister.
	// 其实可以把 Informer 看作是 watch 操作.
	nodeInformer := factory.Core().V1().Nodes()
	informer := nodeInformer.Informer()
```

```

```


关于 scheme

```go
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/apimachinery/pkg/runtime" // 这个包中也有一个 Scheme 成员
```

`runtime.Scheme`与`scheme.Scheme`可以相互赋值.

