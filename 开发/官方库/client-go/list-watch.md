参考文章

1. [理解 K8S 的设计精髓之 List-Watch机制和Informer模块](https://www.jianshu.com/p/234d27d5c1c1)
    - 关于解释`Informer`机制的必要性与优势的前言部分写得很好
    - `Informer`内部运行机制, 十分有条理.
    - 值得阅读
2. [kubernetes 中 informer 的使用](https://www.jianshu.com/p/1e2e686fe363)
    - `kubectl`, `k8s REST API`, `client-go`(`ClientSet`, `Dynamic Client`, `REST`三种方式)等多种方式访问kuber集群获取资源
    - `kubectl get pod -v=9`
    - `Informer`示例demo, 添加注释的代码可见`informer.go.1`文件.
    - `Informer`各组件的作用: `Reflector`, `DeltaIFIFO`, `LocalStore`, `WorkQueue`

kuber: 1.16.2

可以说, 这个示例就是最简的 controller 了, 不需要创建 CRD 资源和对应的 golang 对象(包含 GVK, Spec, Status等), 只监听 kuber 内置的资源对象变动.

`informer`存在的意义就是, 在kuber集群各组件在与apiserver进行通信时添加一个中间缓存层, 减轻apiserver的负载压力. `informer`内置缓存功能, 可保证与apiserver查询到的数据保持一致.

## 执行流程

Etcd存储集群的数据信息, apiserver作为统一入口, 任何对数据的操作都必须经过 apiserver. 客户端(kubelet/scheduler/controller-manager)通过 list-watch 监听 apiserver 中资源(pod/rs/rc等等)的 CURD 事件, 并针对事件类型调用相应的事件处理函数. 

list-watch有两部分组成, 分别是`list`和`watch`. `list`非常好理解, 就是调用资源的list API罗列资源, 基于HTTP短链接实现; `watch`则是调用资源的watch API监听资源变更事件, 基于HTTP长链接实现(http1.1的`chunked`).

`Informer`只会调用kuber中`List`和`Watch`两种类型的API. 

1. Informer在初始化时, 先调用kuber List API 获得某种 resource 的全部Object, 缓存在内存中; 
2. 然后, 调用 Watch API 去watch这种resource, 去维护这份缓存; 
3. 最后, Informer就不再调用kuber的任何API. 

## resync机制

用`List`/`Watch`去维护缓存、保持一致性是非常典型的做法, 但令人费解的是, Informer 只在初始化时调用一次List API, 之后完全依赖 Watch API去维护缓存, 没有任何`resync`机制. 

笔者在阅读Informer代码时候, 对这种做法十分不解. 按照多数人思路, 通过`resync`机制, 重新List一遍 resource下的所有Object, 可以更好的保证 Informer 缓存和 kuber 中数据的一致性. 

咨询过 Google 内部 kube 开发人员之后, 得到的回复是:

> 在 Informer 设计之初, 确实存在一个relist无法去执`resync`操作, 但后来被取消了. 原因是现有的这种 List/Watch 机制, 完全能够保证永远不会漏掉任何事件, 因此完全没有必要再添加relist方法去resync informer的缓存. 这种做法也说明了kuber完全信任etcd. 
