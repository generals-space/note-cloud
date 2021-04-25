# kubebuilder

参考文章

1. [如何看待kubebuilder与Operator Framework(Operator SDK) ？](https://www.zhihu.com/question/290497164)
    - 两者都做了相同的事, sdk中没有 manager 的概念, 只有 controller.
2. [Kubernetes API 编程利器：Operator 和 Operator Framework](https://www.cnblogs.com/yunqishequ/p/12395754.html)
    - 主流的`operator framework`主要有两个: `kubebuilder`和`operator-sdk`
    - 这两者各自的优势
3. [kubebuilder2.0学习笔记——搭建和使用](https://segmentfault.com/a/1190000020338350)
4. [kubebuilder2.0学习笔记——进阶使用](https://segmentfault.com/a/1190000020359577)
    - logr 日志对象
    - eventer事件对象
5. [kubebuilder 使用手册(EN)](https://book.kubebuilder.io/introduction.html)
6. [kubebuilder 使用手册(CN)](https://cloudnative.to/kubebuilder/quick-start.html)
7. [kubebuilder 进阶使用](https://zhuanlan.zhihu.com/p/144978395)

> operator 是描述、部署和管理 kubernetes 应用的一套机制, 从实现上来说, 可以将其理解为 CRD 配合可选的 webhook 与 controller 来实现用户业务逻辑, 即 operator = CRD + webhook + controller

operator 是一种组合概念, 可以说任何使用 CRD 与 controller 实现了某一业务逻辑的工程, 都可以叫做一个 operator.
