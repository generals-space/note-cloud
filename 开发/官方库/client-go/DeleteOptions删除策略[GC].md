# DeleteOptions删除策略[GC]

参考文章

1. [Kubernetes控制器删除策略](https://blog.csdn.net/nangonghen/article/details/102255334)
    - 控制器的删除的3种模式: Foreground, Background, Orphan
    - `curl`命令示例
2. [kubernetes 删除资源对象策略分析](https://blog.csdn.net/qq_21816375/article/details/86500089)
    - 内容同上
3. [Kubernetes之Garbage Collection](https://blog.csdn.net/dkfajsldfsdfsd/article/details/81130786)
    - kubernetes对象并不会产生垃圾, 这里称为"Garbage Collection"不太准确, 实际上它想讲的是在执行删除操作时如何控制对象之间的依赖关系

