# 容器内存的usage_bytes与working_set_bytes

参考文章

1. [cadvisor metrics container_memory_working_set_bytes vs container_memory_usage_bytes](https://blog.csdn.net/u010918487/article/details/106190764)
    - `usage_bytes` > `working_set_bytes`
    - `container_memory_usage_bytes`包含了 cache, 如 filesystem cache, 当存在 mem pressure 的时候, `cache`能够被回收.
    - container 中进程的 OOM 触发是依据`working_set_bytes`的值.
2. [关于kubernetes中监控pod内存的问题](https://blog.csdn.net/qq_34857250/article/details/90378042)
    - 对内存相关的指标给出了比较官方的解释, 以及ta们之间的**大小排序**
3. [kubernetes上报Pod已用内存不准问题分析](https://cloud.tencent.com/developer/article/1637682)
    - `cgroup`与`docker`协作, 精确验证内存占用的场景.
    - 探索过程值得一读
4. [How much is too much? The Linux OOMKiller and “used” memory](https://medium.com/faun/how-much-is-too-much-the-linux-oomkiller-and-used-memory-d32186f29c9d)
    - oom killer 也是根据 container_memory_working_set_bytes 来决定是否oom kill的. 
    - 很有意思的一个示例

## 引言

在使用 prometheus + grafana 进行监控时, 对于内存的监控一般有两个比较常用的指标: `container_memory_working_set_bytes`和`container_memory_usage_bytes`

一般来说, `usage_bytes`都会比`working_set_bytes`要大.

`container_memory_usage_bytes`包含了 cache, 如 filesystem cache, 当存在 mem pressure 的时候, `cache`能够被回收.

`container_memory_working_set_bytes`更能体现出 mem usage, oom killer 也是根据 container_memory_working_set_bytes 来决定是否oom kill的. 

container_memory_max_usage_bytes(最大可用内存) >
container_memory_usage_bytes(已经申请的内存+工作集使用的内存) >
container_memory_working_set_bytes(工作集内存) >
container_memory_rss(常驻内存集)

## 记录

k8s 的容器资源指标是通过 kubelet -> cAdvisor -> runc/libcontainer -> cgroup 得到的, 即, 最终来源就是从 cgroup 目录下查到的.

进到`/sys/fs/cgroup`目录, 该目录中存放着各种资源子系统如 cpu, memory, pids 等.

每种子系统目录下都存在着占用此类资源的分组, 一般会同时存在`docker`和`kubepods`, 不过由于部署了 kuber 之后, 就不会再在宿主机上直接通过 docker 创建容器了, 所以`docker`分组下的内容一般为空(这里的空是指没有用 dockerID 命名的目录), 而`kubepods`分组下就会存在当前宿主机上正在运行的 pod 的目录(这些目录的名称中包含了对应的 docker 容器的名称), 而每个 pod 目录下还会存在与之相关的容器的子目录(pause, 和 containers 中包含的容器).

每一个 pod 的数据目录下, `memory.usage_in_bytes`的统计数据是包含了所有的file cache的, total_active_file和total_inactive_file都属于file cache的一部分, 并且这两个数据并不是业务真正占用的内存, 只是系统为了提高业务的访问IO的效率, 将读写过的文件缓存在内存中, file cache并不会随着进程退出而释放, 只会当容器销毁或者系统内存不足时才会由系统自动回收.

kubectl cAdvisor的`container_memory_usage_bytes` = `memory.usage_in_bytes` - `memory.stat(total_inactive_file)`.

> `memory.stat`文件中包含很多数据, `total_inactive_file`只是其中一个.

但是只减去`total_inactive_file`的内存也是不够的, 因为还剩下`total_active_file`没有减去, 所以`container_memory_usage_bytes`反应 pod 真实的内存占用是不准确的.

## 参考文章3的操作分析

参考文章3对于自己的判断进行了认证, 作者通过`cgroup`命令手动创建资源组, 并通过`docker run`的`--cgroup-parent`选项手动指定目标资源组, 创建了一个测试容器.

对一个大文件进行`grep`, `memory.stat(total_inactive_file)`中的值就会增加相应大小的数值, 再`grep`一次, 该文件占用的内存就会转移到`memory.stat(total_inactive_file)`字段.

猜测是第1次`grep`进行了全文读取, kernel 便将该文件的内容全部放在了 cache, 还是`inactive`的 cache. 第2次`grep`, kernel 就将该文件的内容从`inactive`区域转移到了`active`区域, 也是比较容易理解的.

然后他通过在容器中执行`echo 3 > /proc/sys/vm/drop_caches`, 手动触发了一次内存回收, active(file) 和 inactive(file)都被回收了.
