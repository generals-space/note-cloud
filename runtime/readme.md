参考文章

1. [lxcfs容器隔离技术实现原理分析之loadavg、cpuonline](https://blog.csdn.net/ZVAyIVqt0UFji/article/details/103193083)
2. [使用 Go 和 Linux Kernel 技术探究容器化原理](https://zhuanlan.zhihu.com/p/512715825)
3. [nestybox/sysbox](https://github.com/nestybox/sysbox)
    - 与 runc、gVisor 同级的 runtime, 可以被 docker/kubernetes 使用.

容器云是以云的概念看容器, 有些对象的称谓是不同的.

如在kubernetes中, 存在pv和pvc, 而在容器云中对应的概念是: 存储卷, 存储声明和存储类.

而且容器云平台上一般不会允许用户直接通过kubectl -> apiserver对集群进行控制.
