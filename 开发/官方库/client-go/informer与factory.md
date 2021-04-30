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
