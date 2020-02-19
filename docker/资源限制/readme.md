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

CGroup: Control Group

`/sys/fs/cgroup`目录应该是安装`libcgroup`工具后出现的, 其下有各种控制器subsystem(cpu, memory, blk等). 每种subsystem下可以存储的是该subsystem所有可选选项, 以及所有设置了此subsystem下cgroup的subsystem名称, 各特定的subsystem下是各该subsystem本身的设置值.
