# builder问题处理

## 1. 

```console
$ kubebuilder init 
2020/06/09 09:50:20 failed to initialize project: error finding current repository: could not determine repository path from module data, package data, or by initializing a module: go: cannot determine module path for source directory /home/project (outside GOPATH, no import comments)
```

在一个空目录中执行`kubebuilder init`失败, 出现上述报错.

这个命令需要 go module 模块支持, 且目标工程的 go.mod 文件必须事件存在, 使用`go mod init XXX`初始化完成后, 再执行`kubebuilder init`就可以了.

## 2. init 失败: module requires Go 1.13

```
$ kubebuilder init
Writing scaffold for you to edit...
Get controller runtime:
$ go get sigs.k8s.io/controller-runtime@v0.5.0
# sigs.k8s.io/controller-runtime/pkg/client/apiutil
/usr/local/gopath/pkg/mod/sigs.k8s.io/controller-runtime@v0.5.0/pkg/client/apiutil/dynamicrestmapper.go:48:5: undefined: errors.As
/usr/local/gopath/pkg/mod/sigs.k8s.io/controller-runtime@v0.5.0/pkg/client/apiutil/dynamicrestmapper.go:185:17: undefined: errors.As
/usr/local/gopath/pkg/mod/sigs.k8s.io/controller-runtime@v0.5.0/pkg/client/apiutil/dynamicrestmapper.go:196:16: undefined: errors.As
note: module requires Go 1.13
2020/06/09 10:36:46 failed to initialize project: exit status 2
```

`go mod init`生成的 go.mod 文件中, go 的版本为`1.12`, 但是在`kubebuilder init`时出现上述错误.

看起来像是`kubebuilder`要求的 go 版本不一致. 将本地 go 版本升级至 1.13 再 init 即可.
