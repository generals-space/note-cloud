# OOM相关

参考文章

1. [Kubernetes因限制内存配置引发的错误](https://cloud.tencent.com/developer/article/1411527)
    - 容器 oom 的两种场景
2. [Kubernetes 内存资源限制实战](https://www.cnblogs.com/xingzheanan/p/14837165.html)
3. [Kubernetes 触发 OOMKilled(内存杀手)如何排除故障](https://cloud.tencent.com/developer/article/2314583)
    - 当 Kubernetes 集群中的容器超出其内存限制时，Kubernetes 系统可能会终止该容器，并显示"OOMKilled"错误
    - kernel: Out of memory: Kill process 进程号 (bigmem) score 分值 or sacrifice child
    - kernel: Memory cgroup out of memory: Kill process
    - `/proc/sysrq-trigger`可以在没有发生OOM的时候手动触发, 至少会有一个进程被杀死.

