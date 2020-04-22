# crd开发

参考文章

1. [Kubernetes Deep Dive: Code Generation for CustomResources](https://blog.openshift.com/kubernetes-deep-dive-code-generation-customresources/)
    - `code-generator`项目的`readme.md`中提到的参考文档, 不过不建议用于入门.
2. [k8s代码自动生成过程的解析(Code-Generator)](http://blog.xbblfz.site/2018/09/19/k8s%E4%BB%A3%E7%A0%81%E8%87%AA%E5%8A%A8%E7%94%9F%E6%88%90%E8%BF%87%E7%A8%8B%E7%9A%84%E8%A7%A3%E6%9E%90/)
    - 参考文章1的中文版, 可以稍微读一下.
3. [Kubernetes CRD 系列：Client-Go 的使用](https://liqiang.io/post/kubernetes-all-about-crd-part03-usage-for-client-go-d831d52e#CRD%20%E5%A6%82%E4%BD%95%E4%BD%BF%E7%94%A8%20Typed%20Client)
    - 讲解了`Dynamic Client`和`Typed Clients`两种client对象, 及创建CRD工程时使用`Typed Clients`的原因.
    - 此作者的各文章都非常有深度, 值得一看(不过结构, 内容什么的不太容易理解, 只能作为一个引路人)
4. [使用 code-generator 为 CustomResources 生成代码](https://blog.tianfeiyu.com/2019/08/06/code_generator/)
    - CRD工程的创建步骤及实例.
    - 介绍了`code-generator`中各生成器的作用及使用场景.
    - 执行`generate-groups.sh`生成代码时的4个参数的作用.
5. [Kubernetes CRD Operator 实现指南](https://zhuanlan.zhihu.com/p/38372448)
    - 诸多与kuber同级的编排工具: Mesos(两层式调度器架构), Sparrow(去中心化架构), Hawk(混合式调度器架构), Nomad(共享状态的调度器架构), 及ta们各自的优缺点.
    - CRD的默认规则: name 通常是 plural 和 group 的结合; 另外, 一般来说 CRD 的作用域是 namespaced 就可以了; 还有 kind 一般采用驼峰命名法等..
6. [kubernetes/sample-apiserver工程readme - When using go 1.11 modules](https://github.com/kubernetes/sample-apiserver/#when-using-go-111-modules)
    - `sample-apiserver` v1.17+
    - 使用`code-generator`时无法使用`go mod`创建工程
    - `Note, however, that if you intend to generate code then you will also need the code-generator repo to exist in an old-style location. One easy way to do this is to use the command go mod vendor to create and populate the vendor directory.`
7. [k8s自定义controller三部曲之一:创建CRD（Custom Resource Definition）](https://blog.csdn.net/boling_cavalry/article/details/88917818)
    - 示例工程不错, 给出了`code-generator`生成代码后, 额外的`main.go`, `controller.go`和`signal.go`等文件的编写方法.

如下是`sample-controller`中的crd部署文件.

```yaml
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  ## name=`spec.names.plural+spec.group`
  name: foos.samplecontroller.k8s.io
spec:
  ## CRD工程中, CRD资源可能不只一种, 但group应当是同一个.
  group: samplecontroller.k8s.io
  version: v1alpha1
  names:
    kind: Foo
    ## plugral(复数)
    plural: foos
  scope: Namespaced
```

## 准备步骤

我们首先要创建一个CRD工程目录, 这里假设CRD名为`PodGroup`. 需要注意的是, 貌似使用`code-generator`的项目只能位于`$GOPATH/src`, 所以目前只能先这么做, 见参考文章5.

```console
$ mkdir $GOPATH/src/podgroup
```

按照参考文章3和4中所说, 确定`group`与`version`两个变量的值. 理论上来说, 这两个值可以是随机取的. 这里将`group`定为`testgroup.k8s.io`, `version`定为`v1`. 

> 还记得RBAC中`Role/ClusterRole`的权限列表中有`apiGroups`字段么? 没错, `apiGroups`就是这里的`group`字段, 每个`group`下拥有一种或多种资源.

然后预创建3个文件, 这3个文件在同一个目录下, 路径为`$GOPATH/src/podgroup/pkg/apis/${group}/${version}/`. 

这里的`group`与上面的`group`不是同一个, 上面的`group`为`CustomResourceDefinition` yaml部署文件中的`spec.group`的全名, 而这里路径中的`group`, 是去除后面的组织名称部分, 即`testgroup`. 之后我们将这个变量称为`groupShortName`, 其实ta也可以是`spec.names.shortNames`数组中的其中一个.

> kuber官方声明的CRD的`group`名称, 一般以`k8s.io`结尾, coreos声明的CRD的`group`都以`coreos.com`结尾.

```console
$ mkdir -p $GOPATH/src/podgroup/pkg/apis/testgroup/v1/
$ touch $GOPATH/src/podgroup/pkg/apis/testgroup/v1/doc.go
$ touch $GOPATH/src/podgroup/pkg/apis/testgroup/v1/register.go
$ touch $GOPATH/src/podgroup/pkg/apis/testgroup/v1/types.go
```

> 之所以路径中填写的是`groupShortName`, 是因为一个CRD工程中可能不只有一种自定义资源, `types.go`文件中可以写很多, 所以放在以`groupShortName`为名的目录下, 才比较合理.

文件的内容就不在这里贴了, 可以见`PodGroup`项目的`pkg/apis`目录. 除了具体的类型定义, 其他的像`import`的内容, `+genclient`和`+k8s`这种编译标记, 都拷贝过去.

## 代码生成

需要注意的是`code-generator`的代码生成步骤. 

网上好像都没有文章明确地讲过具体步骤, 有的说需要把`code-generator`目录下的`vendor`和`hack`目录拷贝到CRD工程目录, 然后执行`vendor/k8s.io/code-generator/generate-groups.sh xxx`(但是`code-generator`的vendor根本没有ta本身的工程); 参考文章3根本就没说, 直接在当前目录就执行了`./generate-groups.sh xxx`, 这意思是先拷贝一份`code-generator`工程?

首先, `generate-groups.sh`要求`code-generator`和`apimachinery`两个工程在`$GOPATH/src/k8s.io/`目录下(我们的CRD工程也要放在`$GOPATH`目录下.), `go mod`形式的依赖管理无效. 否则在执行脚本时会出现`Hit an unsupported type invalid type for invalid type`的问题.

然后执行如下命令

```console
$ $GOPATH/src/k8s.io/code-generator/generate-groups.sh all podgroup/pkg/client podgroup/pkg/apis testgroup:v1
```

执行命令时所在的目录没有强制要求, 另外, `testgroup:v1`, 对应了`podgroup`项目中的`pkg/apis/testgroup/v1`目录, 如果两者不一致, 会报`Error: Failed making a parser: unable to add directory "podgroup/xxx": unable to import "podgroup/xxx": cannot find package "podgroup/xxx"`.

然后可正常生成代码.

```console
$ $GOPATH/src/k8s.io/code-generator/generate-groups.sh all podgroup/pkg/client podgroup/pkg/apis testgroup:v1
Generating deepcopy funcs
Generating clientset for testgroup:v1 at podgroup/pkg/client/clientset
Generating listers for testgroup:v1 at podgroup/pkg/client/listers
Generating informers for testgroup:v1 at podgroup/pkg/client/informers
```

## 完善补充

生成代码后(生成的代码树结构网上一大堆), 需要补充其他文件.

首先是先把工程目录从`GOPATH`下移出, 并创建`go.mod`文件进行`go modules`初始化.

然后创建`pkg/signals/signal.go`文件.

再创建根目录下的`main.go`和`controller.go`.

另外, 代码生成完成后, 对`PodGroup{}`及`PodGroupList{}`的成员进行修改就不再需要重新生成了.
