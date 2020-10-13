参考文章

1. [kubernetes/code-generator](https://github.com/kubernetes/code-generator)
2. [kubernetes/sample-controller](https://github.com/kubernetes/sample-controller)
3. [kubernetes/sample-apiserver](https://github.com/kubernetes/sample-apiserver)

`sample-controller`与`sample-apiserver`的 readme 中, 都存在一个`When using go 1.11 modules`小节, 内容如下

```
Note, however, that if you intend to generate code then you will also need the code-generator repo to exist in an old-style location. One easy way to do this is to use the command go mod vendor to create and populate the vendor directory.
```
