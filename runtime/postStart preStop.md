参考文章

1. [Pod的Terminated过程](https://www.cnblogs.com/orchidzjl/p/11791883.html)
    - sigkill是不能被捕获的, 程序收到这个信号后, 一定会退出
2. [sigma敏捷版系列文章：kubernetes grace period 失效问题排查](https://developer.aliyun.com/article/609813)
    - `PostStart`执行的时机是在容器启动以后，但是并不是等容器启动完成再执行(可以说`PostStart`和`Entrypoint`是并行执行的)

Kubernetes 在创建容器后立即发送`postStart`事件. 但是, **不能保证`postStart`处理程序在容器的`entrypoint`调用之前被调用**. 相对于容器的代码, `postStart`处理程序以**异步方式**运行, 但 Kubernetes 对容器的管理会阻塞直到`postStart`处理程序完成. 容器的状态直到`postStart`处理程序完成后才会设置为`RUNNING`. 

Kubernetes 在容器终止之前立即发送`preStop`事件. Kubernetes 对容器的管理一直阻塞直到`preStop`处理程序完成, 除非 Pod 的宽限期过期. 
