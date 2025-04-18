[Linux:如何限制进程的CPU使用率?](https://www.techforgeek.info/how_to_limit_cpu_usage.html)
[Linux的公平调度（CFS）原理 - kummer话你知](https://www.jianshu.com/p/673c9e4817a8)
    - linux调度算法的演进过程: O(n) -> O(1) -> CFS
    - CFS的运行原理 - 红黑树
2. [docker cgroup技术之cpu和cpuset](https://www.cnblogs.com/charlieroro/p/10281469.html)
    - 介绍了cgroup中的cpu, cpuset两个子系统各自的作用及两个子系统中比较重要的几个参数的含义和使用方法
1. [DOCKER基础技术：LINUX CGROUP](https://coolshell.cn/articles/17049.html)
    - cgroup文件系统的挂载及使用方法
    - C语言示例演示CPU, 内存占用
[浅谈Linux Cgroups机制](https://zhuanlan.zhihu.com/p/81668069)
    - 比较全, 列举比较多, 但对概念解释不清晰, 应该是没实践过.
[深入理解 Linux Cgroup 系列（一）：基本概念](https://www.cnblogs.com/ryanyangcs/p/11198140.html)
    - 比较易懂
[Tencent/TencentOS-kernel](https://github.com/Tencent/TencentOS-kernel)
    - readme里讲了一些高大上的东西.
[获取 Docker container 中的资源使用情况](https://zhuanlan.zhihu.com/p/35914450)
[Java 10发布后，Docker容器管理能力得到显著增强](http://www.dockone.io/article/5931)
[kubernetes的资源管理](http://bazingafeng.com/2017/12/04/the-management-of-resource-in-kubernetes/)

[docker资源限制使用和性能压测](https://zhuanlan.zhihu.com/p/641334756)

CGroup: Control Group

`/sys/fs/cgroup`目录应该是安装`libcgroup`工具后出现的, 其下有各种控制器subsystem(cpu, memory, blk等). 

每种subsystem下可以存储的是该subsystem所有可选选项(cpu下的`cfs_period_us`, `shares`这种), 以及所有设置了此subsystem下cgroup的subsystem名称, 各特定的subsystem下是各该subsystem本身的设置值.

进到`/sys/fs/cgroup`目录, 该目录中存放着各种资源子系统如 cpu, memory, pids 等.

每种子系统目录下都存在着占用此类资源的分组, 一般会同时存在`docker`和`kubepods`, 不过由于部署了 kuber 之后, 就不会再在宿主机上直接通过 docker 创建容器了, 所以`docker`分组下的内容一般为空(这里的空是指没有用 dockerID 命名的目录), 而`kubepods`分组下就会存在当前宿主机上正在运行的 pod 的目录(这些目录的名称中包含了对应的 docker 容器的名称), 而每个 pod 目录下还会存在与之相关的容器的子目录(pause, 和 containers 中包含的容器).

----

`docker`, `kubelet`的`cgroup driver`的可选值都为`systemd`与`cgroupfs`.

`dockerd`使用的 cgroup driver 默认为`cgroupfs`.

