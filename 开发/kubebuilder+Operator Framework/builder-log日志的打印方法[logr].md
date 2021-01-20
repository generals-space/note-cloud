# builder-logger日志对象的使用方法

参考文章

1. [kubebuilder2.0学习笔记——进阶使用](https://segmentfault.com/a/1190000020359577)
    - kubebuilder内置的`github.com/go-logr/logr`日志对象的使用方法

使用 kubebuilder 初始构建的工程中, `XXX_controller.go`文件的`Reconcile()`方法内容如下.

```go
func (r *MyClusterReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	_ = context.Background()
	_ = r.Log.WithValues("mycluster", req.NamespacedName)

	// your logic here

	return ctrl.Result{}, nil
}
```

其中`r.Log.WithValues()`可以得到日志对象, 我们需要使用这个对象打印日志信息

```go
func (r *MyClusterReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	_ := context.Background()
	logger := r.Log.WithValues("mycluster", req.NamespacedName)

    // your logic here
    // 日志格式为
    // object info	{"mycluster": "kube-system/mycluster-sample", "info": {"namespace": "kube-system", "name": "mycluster-sample"}}
    logger.Info("object info", "info", req.NamespacedName)

	return ctrl.Result{}, nil
}
```

需要了解的是, `logger.Info()`接受的参数中, 第1个参数为前缀, 之后的参数都是 k-v 对, 日志中会像 json 中的键值对一样打印出来.
