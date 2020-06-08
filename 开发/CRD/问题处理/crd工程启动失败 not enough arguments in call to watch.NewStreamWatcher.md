# crd工程启动失败 not enough arguments in call to watch.NewStreamWatcher

参考文章

1. [watch.NewStreamer error](https://github.com/kubernetes/client-go/issues/584)

```console
$ go run main.go controller.go 
# k8s.io/client-go/rest
/usr/local/gopath/pkg/mod/k8s.io/client-go@v11.0.0+incompatible/rest/request.go:598:31: not enough arguments in call to watch.NewStreamWatcher
	have (*versioned.Decoder)
	want (watch.Decoder, watch.Reporter)
```

估计是go.mod文件中声明的依赖版本不匹配(比如`incompatible`标记的记录), 很大可能是`client-go`的版本, 可以查查.
