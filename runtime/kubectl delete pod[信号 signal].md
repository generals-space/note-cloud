参考文章

1. [Pod的Terminated过程](https://www.cnblogs.com/orchidzjl/p/11791883.html)
    - shell 程序不转发 signals, 也不响应退出信号SIGTERM和SIGKILL
    - kernel会为每个进程(除了init进程)加上默认的 signal handler, 所以kill是可以杀死一个 shell 脚本进程的
    - SIGINT SIGTERM SIGKILL区别
    - SIGKILL 是不能被捕获的, 程序收到这个信号后, 一定会退出
    - k8s中Pod的终止过程, postStart 与 preStop 的处理时机
2. [Docker 里的进程为什么没有处理 TERM 信号](https://jan365.org/post/process-in-docker-does-not-handle-term/)
3. [优雅地关闭kubernetes中的nginx](https://segmentfault.com/a/1190000008233992)
    - 因为SIGKILL信号是直接发往系统内核的, 应用程序没有机会去处理它(其实也没有机会接收到)
    - 当优雅退出时间超时了, 任何pod中正在运行的进程会被发送`SIGKILL`信号被杀死。
4. [sigma敏捷版系列文章：kubernetes grace period 失效问题排查](https://developer.aliyun.com/article/609813)
    - Container Lifecycle Hooks 的详细介绍, postStart, preStop 的调用时机
    - 删除 Pod 时的详细的过程

`kubectl delete pod`就是对pod中启动的 pid 1 进程发送 TERM 信号, 30s后无响应则发 KILL 信号强制退出.


我们知道, 在`dockerfile`中使用`cmd/entrypoint` 中指定一个脚本, 然后由该脚本启动一个服务时, pid 为 1 的进程将是那个脚本.

> 注意: 虽然`cmd/entrypoint`是通过`sh -c 'xxx'`的方式运行的指令, 但 pid 为 1 的进程的确是 xxx 而不是`sh`(很像`fork`, 把原来的进程直接替换掉了)

虽然可以指定一个`xxx.sh`脚本作为 pid 为 1 的进程, 但是一般脚本里也不会用`trap`去捕获信号, 所以`kubectl delete`这种类型的Pod的时候一般会超过30s, 有机会可以使用如下命令试验一下.

```
time k delete pod xxx
```

但是这种 shell 脚本形式的 init 进程对 kill 信号其实是完全无反应的, 无论是发送 SIGTERM, SIGINT, 还是 SIGKILL 信号都不会影响到ta.

那么 kubelet 是怎么终止这样的容器呢?

按照参考文章3中所说, 当向 pid 1 进程发送 SIGTERM 信号过了30s后, pod 仍旧没有停止, kubelet 就会向 pod 中所有进程发送 SIGKILL 信号.

不过由于 SIGKILL 信号没有办法被进程捕获到, 所以没有办法做个实验验证一下. 只能有机会去查看 runc 底层的代码了.

