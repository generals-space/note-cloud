参考文章

1. [理解 K8S 的设计精髓之 List-Watch机制和Informer模块](https://www.jianshu.com/p/234d27d5c1c1)
    - 关于解释`Informer`机制的必要性与优势的前言部分写得很好
    - `Informer`内部运行机制, 十分有条理.
    - 值得阅读

reflector 是 client-go 的核心.

普通 operator 是与 apiserver 建立通信的, 那么 apiserver 自己也是要用 client-go 库的, 那 apiserver 是怎么用的?

如果一个 reflector 对象, 只能针对单个类型的资源. 但是查看 Reflector{} 的结构, 会发现ta的是一个通用结构.

reflector 是 client-go 的核心, tools/cache/reflector.go -> Reflector.ListAndWatch() 启动后会先全量同步, 然后 watch, 都是调用的自身的 list-watcher 成员.

如果是普通的 operator, 那么 list-watcher 是通过 rest-client 实现的, 而如果是 apiserver, 则是通过 etcd storage 实现的. 

可见 staging/src/k8s.io/apiserver/pkg/storage/cacher/cacher.go -> cacherListerWatcher{} 的 List/Watch 方法实现.

reflector 与 queue 不直接相关.
