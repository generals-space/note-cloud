# go mod 引入 kubernetes 作为依赖库

<!key!>: {d03f75fd-1b5d-46f2-a100-85566ec39d51}
<!link!>: {71d63267-269b-4154-92fc-daad8767aae4}

参考文章

1. [k8s.io/api@v0.0.0: reading k8s.io/api/go.mod at revision v0.0.0: unknown revision v0.0.0](https://github.com/kubernetes/kubernetes/issues/90358)
2. ['unknown revision v0.0.0' errors, seemingly due to 'require k8s.io/foo v0.0.0'](https://github.com/kubernetes/kubernetes/issues/79384)

kuber: v1.17.3

最近在学着写扩展调度器, 找到一个示例工程[everpeace/k8s-scheduler-extender-example](https://github.com/everpeace/k8s-scheduler-extender-example), 其中需要引入`k8s.io/kubernetes/pkg/scheduler/apis/extender/v1`包.

于是我在命令行使用`go get`尝试下载 kuber 工程, 但是却报如下错误.

```console
$ go get -v k8s.io/kubernetes@v1.17.3
go: finding k8s.io v1.17.3
go get: k8s.io/kubernetes@v1.17.3 requires
	k8s.io/api@v0.0.0: reading https://goproxy.io/k8s.io/api/@v/v0.0.0.mod: 404 Not Found
```

我初始以为是 goproxy 镜像设置得不合适, 当前镜像站没有保留相应的 tag 导致404, 尝试使用 direct 然后开代理下载, 结果还是不行.

后来在网上找到了参考文章1. 

其实我之前也遇到过这个问题, 可以见本文的链接文章. 当时找到的关键是下面这句话

```
k8s.io/kubernetes is not primarily intended to be consumed as a module. Only the published subcomponents are (and go get works properly with those).
```

官方kubernetes仓库并不能直接作为module使用, 所以无法通过`go get`或是`go mod download`下载其中的package, 只有某些已发布的子组件可以(如`k8s.io/api`, `k8s.io/client-go`).

但是之前的文章并没有解决这个问题, 好在参考文章1中有人给出了各详细的解释.

```
This is caused by depending on k8s.io/kubernetes directly as a library, which is not supported. The components intended to be used as libraries are published as standalone modules like k8s.io/api, k8s.io/apimachinery, k8s.io/client-go, etc, and can be referenced directly.
```

同时又给出了参考文章2的链接.

参考文章2中指出, kubernetes 仓库本身将很多依赖指向了本地

```
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

这样, ta是不在乎这些依赖包的版本的.

但是当我们在`require k8s.io/kubernetes v1.17.3`时, ta想下载自己的`k8s.io/api`依赖, 但是ta本地的`./staging/src/k8s.io/api`又没有版本号, 根本不是一个合法的 golang 仓库, 因为就出现了我所遇到的问题.

为了解决这个问题, 参考文章2指出, 可以使用`replace`指令直接指定`k8s.io/api`与 kubernetes 相应的版本, 让 go mod 在下载 kubernetes 时使用我们指定的`k8s.io/api`版本.

```go
require (
	k8s.io/kubernetes v1.17.3
)
replace k8s.io/api => k8s.io/api v0.17.3
```

但是紧接着又出现了如下问题

```
$ go mod download
go: k8s.io/kubernetes@v1.17.3 requires
	k8s.io/apiextensions-apiserver@v0.0.0: reading k8s.io/apiextensions-apiserver/go.mod at revision v0.0.0: unknown revision v0.0.0
```

所以我们需要把所有 kubernetes 内置的依赖库都`replace`一遍.

```

```
