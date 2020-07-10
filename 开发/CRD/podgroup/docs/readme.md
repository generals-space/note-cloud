按照后面两篇文章, 准备好 code-generator 和 apimachinery 两个工程, 然后重新生成代码, 就可以执行了.

## FAQ

```
$ go run main.go
build command-line-arguments: cannot load podgroup/pkg/client/clientset/versioned: malformed module path "podgroup/pkg/client/clientset/versioned": missing dot in first path element
```

这是因为还没生成代码, 或是生成了代码但不在当前工程目录下, 到`$GOPATH`下去找找, 然后拷贝到自己的工程中.
