# builder使用方法

参考文章

1. [Kubernetes API 编程利器：Operator 和 Operator Framework](https://www.cnblogs.com/yunqishequ/p/12395754.html)
    - 主流的 operator framework 主要有两个：kubebuilder 和 operator-sdk
    - code generator 依赖注释生成代码, 并给出了示例中使用的注释
2. [利用 kubebuilder 优化 Kubernetes Operator 开发体验](https://caicloud.io/blog/5d02213311f1d9002c8543ef)
    - 才云科技
    - ~~kubebuilder 是开发 Operator 的框架~~ 准确来说, kubebuilder 是开发 Controller Manager 的框架, CM 会管理一个或者多个 Operator
3. [利用 kubebuilder 优化 K8s Operator 开发体验](http://k8smeetup.com/article/VkpAij8gP)
    - 参考文章2的转载文章(有图片)
4. [kubebuilder(1)-安装和使用](https://www.jianshu.com/p/de4dd9c9ad47)
    - 分析了在 init, create api 不同阶段代码的变化, 分步介绍.
    - 解释了 webhook, crt, finalizer 等各组件的作用:
        1. webhook: 当添加和修改一个Object时, 需要对Object的合法性进行判断, 可以通过webhook的framework来进行合法性的判定; 
        2. crt: 用于解决webhook访问k8s时所需要的证书问题, 官网也建议使用crt-manager解决证书问题; 
        3. finalizer: 由于这个Object可能创建一些其他的resource(比如pod), 在删除之前, 需要做一些清理工作, finalizer就是实现这个清理的framework代码; 
5. [使用 Kubebuilder 构建 Kubernetes CRD Controller](https://blog.ihypo.net/15645917310391.html)
    - 给出了一个完成的 kubebuilder 完成的工程示例
6. [kubebuilder2.0学习笔记——搭建和使用](https://segmentfault.com/a/1190000020338350)
7. [kubebuilder2.0学习笔记——进阶使用](https://segmentfault.com/a/1190000020359577)
8. [kubebuilder 官方文档](https://book.kubebuilder.io/introduction.html)

`kubebuilder`与`operator-sdk`都是二进制文件, 可以不下载源码直接执行. 另外`kubebuilder`的`Makefile`中调用了`kustomize`, 所以也需要下载`kustomize`的可执行文件.

kubebuilder: 2.3.1
kustomize: 3.6.1

```console
$ kubebuilder version
Version: version.Version{KubeBuilderVersion:"2.3.1", KubernetesVendor:"1.16.4", GitCommit:"8b53abeb4280186e494b726edf8f54ca7aa64a49", BuildDate:"2020-03-26T16:42:00Z", GoOs:"unknown", GoArch:"unknown"}
```

`kubebuilder version`会打印出`KubernetesVendor`字段, 最初以为目标集群版本需要与这个值匹配, 但是在 1.17.2 集群中部署我们的工程, 并没有出错. 要么, 可能就是这表示的是最低版本, 低于这个版本才会报错...???

## 1. init

```
kubebuilder init
kubebuilder init --domain=generals.space
```

这个命令会初始化当前目录所在工程, 需要该目录拥有 go.mod 文件, 否则会报错.

`--domain`: 不难理解, 主要就是一个标识的作用, 没有实际意义. 默认值为"my.domain"

> `kubebuilder`可以指定要创建 CRD 的 domain, group 信息, 都有默认值, 不会根据 go.mod 或是当前工程目录名称来取名字...

## 2. create api

```
kubebuilder create api --group apps --version v1alpha1 --kind PodGroup
```

这次没有默认值了, GVK三者都是必填项. 参考文章1有都这3者的含义有明确解释.

这一步会在当前目录生成`api`目录, 其下有`types.go`, `zz_generated.deepcopy`等文件. 同时也会生成`controllers`目录, 其下会包含`podgroup_controller.go`文件及对应的单元测试文件.

## 3. 编写代码

`types.go`仍然存放着 CRD 的结构声明. 我们要关注的, 目前只有`podgroup_controller.go`这个文件, ta 与传统使用 code-generator 生成的代码而编写的 controller 不同.

## 4. 运行&部署

到了这一步, 网上的文章说的几乎都是`make && make install`, 但我们这里将其中的步骤拆分开来.

首先当然是`make`所表示的, build 我们的代码, 其中也包含了编译与验证的过程.

然后`make install`调用`kustomize`, 生成 CRD 的 yaml 部署文件, 同时调用`kubectl apply -f`直接将其部署到集群中.

我们首先生成 CRD 的部署文件(这一步最好看一下`Makefile`中的做法, 有一步`manifests`必不可少), 再运行我们构建好的可执行文件(也可以执行`go run main.go`), 最后创建 CR 实例. 在`config/sample`目录下会有一个 CRD 实例的部署文件, 不过貌似其结构在`create api`的时候就已经固定了, 之后`make`与`make install`并不会自动修改其中的内容.

