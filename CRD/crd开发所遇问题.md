# crd开发所遇问题

```
$ /usr/local/gopath/src/k8s.io/code-generator/generate-groups.sh all ./pkg/client ./pkg/apis testgroup.k8s.io:v1
Generating deepcopy funcs
F1223 02:31:46.140237   86710 deepcopy.go:885] Hit an unsupported type invalid type for invalid type, from ./pkg/apis/testgroup.k8s.io/v1.PodGroup

[code-generator: deepcopy-gen invalid type for invalid type](https://github.com/kubernetes/kubernetes/issues/79149)
[Generator only works if k8s.io/apimachinery is in the $GOPATH](https://github.com/kubernetes/code-generator/issues/21)
```

参考文章

1. [code-generator: deepcopy-gen invalid type for invalid type](https://github.com/kubernetes/kubernetes/issues/79149)
2. [Generator only works if k8s.io/apimachinery is in the $GOPATH](https://github.com/kubernetes/code-generator/issues/21)
    - 参考文章1提到本文, 且确认有效.

## 

```console
$ go run main.go controller.go 
# k8s.io/client-go/rest
/usr/local/gopath/pkg/mod/k8s.io/client-go@v11.0.0+incompatible/rest/request.go:598:31: not enough arguments in call to watch.NewStreamWatcher
	have (*versioned.Decoder)
	want (watch.Decoder, watch.Reporter)
```

参考文章

1. [watch.NewStreamer error](https://github.com/kubernetes/client-go/issues/584)

估计是go.mod文件中声明的依赖版本不匹配(比如`incompatible`标记的记录), 很大可能是`client-go`的版本, 可以查查.
