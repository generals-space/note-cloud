# code-generator 代码生成失败 Hit an unsupported type invalid type for invalid type

参考文章

1. [code-generator: deepcopy-gen invalid type for invalid type](https://github.com/kubernetes/kubernetes/issues/79149)
2. [Generator only works if k8s.io/apimachinery is in the $GOPATH](https://github.com/kubernetes/code-generator/issues/21)
    - 参考文章1提到本文, 且确认有效.

问题描述

```
$ /usr/local/gopath/src/k8s.io/code-generator/generate-groups.sh all ./pkg/client ./pkg/apis testgroup.k8s.io:v1
Generating deepcopy funcs
F1223 02:31:46.140237   86710 deepcopy.go:885] Hit an unsupported type invalid type for invalid type, from ./pkg/apis/testgroup.k8s.io/v1.PodGroup

[code-generator: deepcopy-gen invalid type for invalid type](https://github.com/kubernetes/kubernetes/issues/79149)
[Generator only works if k8s.io/apimachinery is in the $GOPATH](https://github.com/kubernetes/code-generator/issues/21)
```

- code-generator: v0.17.0
- apimachinery: v0.17.0

`code-generator`要生成代码, 需要`GOPATH`目录下存在`code-generator`和`apimachinery`两个工程, 只在CRD工程的`vendor`目录下不管用.

