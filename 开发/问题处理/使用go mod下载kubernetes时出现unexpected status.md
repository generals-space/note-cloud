# 使用go mod下载kubernetes时出现unexpected status

<!key!>: {71d63267-269b-4154-92fc-daad8767aae4}

<!link!>: {d03f75fd-1b5d-46f2-a100-85566ec39d51}

参考文章

1. [golang官方文档 proxy.golang.org: @v/list includes versions with no corresponding .info file](https://github.com/golang/go/issues/34033)
2. [你好，我设置了GOPROXY为goproxy.cn后，go get 下载gorm出现依赖下载404 #25](https://github.com/goproxy/goproxy.cn/issues/25)
3. [proxy.golang.org: unexpected status 410 Gone](https://blog.csdn.net/shida_csdn/article/details/100056056)
4. [golang官方文档 proxy.golang.org: unexpected status 410 Gone](https://github.com/golang/go/issues/32461)
5. [cmd/go: go get failing with vague error 'unknown revision v0.0.0' with kubernetes](https://github.com/golang/go/issues/32776)
6. ['unknown revision v0.0.0' errors, seemingly due to 'require k8s.io/foo v0.0.0'](https://github.com/kubernetes/kubernetes/issues/79384)


## 场景描述

golang版本: 1.12.14
kuber源码版本: 1.16.0

在阅读proxy组件源码时, 想实验一下pkg目录下`ipvs`模块的功能, 于是创建了一个测试工程, 命名为`test-ipvs`

`go.mod`

```
module test-ipvs

go 1.12
```

`main.go`

```go
package main

import (
	"k8s.io/kubernetes/pkg/proxy/ipvs"
)

func main(){
	var kernelHandler ipvs.KernelHandler
	kernelHandler = ipvs.NewLinuxKernelHandler()

}
```

只有这两个文件.

工程目录下执行`go mod download`, 没反应...于是想手动执行`go get`, 结果...

```log
$ go get -v k8s.io/kubernetes/pkg/proxy/ipvs
go: k8s.io/kube-proxy@v0.0.0: unexpected status (https://goproxy.cn/k8s.io/kube-proxy/@v/v0.0.0.info): 404 Not Found
...
go: error loading module requirements
```

执行结果中出现无数个`404 Not Found`.

此时`go env`为

```
GO111MODULE=auto
GOPROXY=https://goproxy.cn
```

最开始以为可能是因为golang版本的问题, 好像网上大部分文章介绍的都是go1.13, go1.12可能会出现这样的问题. 于是我换了台电脑, `go env`不变, golang版本为`1.13.5`, 结果...

```log
$ go get -v k8s.io/kubernetes/pkg/proxy/ipvs
go get: k8s.io/kubernetes@v1.17.0 requires
	k8s.io/api@v0.0.0: reading https://goproxy.io/k8s.io/api/@v/v0.0.0.mod: 410 Gone
```

没了404, 又来了410...呵呵

## 问题分析

按照参考文章2中所说, 是因为出现404/410的仓库, 之前存在过的tag又被删除, 结果goproxy镜像站因为缓存了实际上已经不存在了的tag, 再去寻找时就出了问题.

如果golang版本为1.12, 可以将GOPROXY修改为`GOPROXY=direct`, 告诉`go mod`机制不用走镜像站, 直接到源站下载.

如果golang版本为1.13, 则GOPROXY的值支持多值列表(逗号分隔), `GOPROXY=https://goproxy.cn,direct`.

------

但是又出现了如下问题.

```log
$ go get -v k8s.io/kubernetes/pkg/proxy/ipvs
go: k8s.io/cluster-bootstrap@v0.0.0: unknown revision v0.0.0
...
go: error loading module requirements
```

多个依赖包出现了`unknown revision`的错误导致失败.

于是找到了参考文章5和6.

首先kubernetes源码中应该不是用`replace`指令做的, 参考文章5中有如下回复.

> k8s.io/kubernetes is not primarily intended to be consumed as a module. Only the published subcomponents are (and go get works properly with those).

官方kubernetes仓库并不能直接作为module使用, 所以无法通过`go get`或是`go mod download`下载其中的package, 只有某些已发布的子组件可以(应该是指`k8s.io/api`, `k8s.io/kube-controller-manager`这种吧).

另外, 参考文章6中也提到, kubernetes的`go.mod`文件中有很多依赖写成这种形式

```go
require(
	k8s.io/apimachinery v0.0.0
	k8s.io/apiserver v0.0.0
	k8s.io/cli-runtime v0.0.0
	k8s.io/client-go v0.0.0
)

replace(
    k8s.io/apimachinery => ./staging/src/k8s.io/apimachinery
	k8s.io/apiserver => ./staging/src/k8s.io/apiserver
	k8s.io/cli-runtime => ./staging/src/k8s.io/cli-runtime
	k8s.io/client-go => ./staging/src/k8s.io/client-go
)
```

可以看到, 这种做法就只是为了把依赖指向本地而已, `vendor`下的这此包也是根据`staging`目录中创建的软链接...

好了, 迷题终于破解了. 就是说, 要使用`k8s.io/kubernetes/pkg/proxy/ipvs`, 只能拷贝这个目录呗...
