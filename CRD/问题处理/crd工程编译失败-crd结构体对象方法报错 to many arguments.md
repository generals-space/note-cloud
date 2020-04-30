# crd工程编译失败-crd结构体对象方法报错 to many arguments

## 问题描述

生成代码完成, 修了一些更新后编译工程, 报错如下

```
pkg/client/clientset/versioned/typed/ipkeeper/v1/staticips.go:134:5: too many arguments in call to c.client.Put().Namespace(c.ns).Resource("staticipses").Name(staticIPs.ObjectMeta.Name).VersionedParams(&opts, scheme.ParameterCodec).Body(staticIPs).Do
	have (context.Context)
	want ()
```

`pkg/client/clientset/versioned/typed/ipkeeper/v1/staticips.go`

```go
// Get takes name of the staticIPs, and returns the corresponding staticIPs object, and an error if there is any.
func (c *staticIPses) Get(ctx context.Context, name string, options metav1.GetOptions) (result *v1.StaticIPs, err error) {
	result = &v1.StaticIPs{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("staticipses").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}
```

其中CRD资源类型为`staticIPs`, 上面的方法中, receiver为其复数形式.

对比了一个`PodGroup`工程, 发现`Get()`方法中多出了`ctx`参数, 且函数体内部调用的`Do()`方法本来也不应该传入`ctx`的.

猜测可能是因为`code-generator`和`apimachinery`版本太新, 于是到`GOPATH`的这两个工程目录下, 使用如下命令切换到之前验证过的`v0.17.0`

```
git checkout -b v0.17.0 v0.17.0
```

在这个版本中, 两个工程已经都使用了`go.mod`作为依赖管理, 其下并没有自带`vendor`, 需要再使用如下命令生成此目录.

```
go mod vendor
```

> 由于是在`GOPATH`目录中执行命令, 可以可能需要考虑`GO111MODULE`变量的影响.

完成后重新生成代码, 就不再报错了.
