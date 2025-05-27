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
    - 同样提到了`container_memory_working_set_bytes`指标计算不准确的问题, 并且比参数文章1更加精确
4. [How much is too much? The Linux OOMKiller and “used” memory](https://medium.com/faun/how-much-is-too-much-the-linux-oomkiller-and-used-memory-d32186f29c9d)
    - oom killer 也是根据 container_memory_working_set_bytes 来决定是否oom kill的. 
    - 很有意思的一个示例

## 引言

在使用 prometheus + grafana 进行监控时, 对于内存的监控一般有几个比较常用的指标: 

container_memory_max_usage_bytes >
container_memory_usage_bytes >
container_memory_working_set_bytes >
container_memory_rss

还有一个含义比较明确的 container_memory_cache

以一个比较典型的内存泄露场景为例, [k8snetworkplumbingwg/multus-cni v4.1.0](https://github.com/k8snetworkplumbingwg/multus-cni/tree/v4.1.0), 每创建一个 Pod, 就会调用一次该 cni 服务, 同时增长一点内存. 当某个 Pod 因为 crash 重复启动, 该服务的内存就会呈现一种比较平滑的上升趋势.

![](https://gitee.com/generals-space/gitimg/raw/master/2025/fe347855776aa5b94edcaa8615d73eca.png)

> memory limit 值为 1G.

> 图片末尾处"usage_bytes", "working_set"与"cache"出现骤降, 是因为使用`echo 1 > /proc/sys/vm/drop_caches`清空了 page cache.

- `container_memory_max_usage_bytes`比较特殊, ta表示的是历史上`container_memory_usage_bytes`指标的最大值, 因此单调递增且没有毛刺(但没有实际意义).
- `container_memory_cache`含义也比较明确, 应该表示 PageCache(还有各种inode cache?), 用于提高各进程的性能表现. 当操作系统出现内存压力的时候, `cache`部分可以被回收.
- ==============================================================================
- `container_memory_rss`表示程序真实的活跃内存(堆栈程序段数据段等)
- `container_memory_working_set_bytes`包含了`container_memory_rss`, 以及部分`container_memory_cache`(被频繁访问的 page cache).
    - 更接近操作系统的"实际可用内存压力"视角
- `container_memory_usage_bytes`应该是`container_memory_working_set_bytes`+`container_memory_cache`, 表示cgroup眼中的实际使用内存.
    - `container_memory_usage_bytes`早早的到达了 limit 上限, 为保证 working_set 的使用, 就开始回收 cache. 因此在`container_memory_usage_bytes`到1G后, `container_memory_working_set_bytes`与`container_memory_cache`呈现相反的趋势.

下图可以作为`container_memory_usage_bytes` = `container_memory_working_set_bytes` + `container_memory_cache`的证明, 两者相差不大.

![](https://gitee.com/generals-space/gitimg/raw/master/2025/888a0f1b235693e3aee9ba072cbeead1.png)

`container_memory_usage_bytes`更能体现出 mem usage, cgroup 的 oom killer 也是根据 container_memory_usage_bytes 来决定是否oom kill的, 在 cache 被回收到极致, 仍然无法满足 working_set 的申请需求时, 就会被 OOM(不过这个论点还没有合适的监控图表支撑).

## 记录

k8s 的容器资源指标是通过 kubelet -> cAdvisor -> runc/libcontainer -> cgroup 得到的, 即, 最终就是从 cgroup 目录下查到的.

进到`/sys/fs/cgroup`目录, 该目录中存放着各种资源子系统如 cpu, memory, pids 等.

每种子系统目录下都存在着占用此类资源的分组, 一般会同时存在`docker`和`kubepods`, 不过由于部署了 kube 之后, 就不会再在宿主机上直接通过 docker 创建容器了, 所以`docker`分组下的内容一般为空(这里的空是指没有用 dockerID 命名的目录), 而`kubepods`分组下就会存在当前宿主机上正在运行的 pod 的目录(这些目录的名称中包含了对应的 docker 容器的名称), 而每个 pod 目录下还会存在与之相关的容器的子目录(pause, 和 containers 中包含的容器).

| metrics指标                          | cgroup信息                                                   |
| :----------------------------------- | :----------------------------------------------------------- |
| `container_memory_rss`               | `memory.stat(rss/total_rss)`                                 |
| `container_memory_usage_bytes`       | `memory.usage_in_bytes`                                      |
| `container_memory_working_set_bytes` | `memory.usage_in_bytes` - `memory.stat(total_inactive_file)` |

> `memory.stat`文件中包含很多数据, `total_inactive_file`只是其中一个.

每一个 pod 的数据目录下, `memory.usage_in_bytes`的统计数据是包含了所有的 file cache 的(内核的计算方式), `total_active_file`和`total_inactive_file`都属于 file cache 的一部分, 并且这两个数据并不是业务真正占用的内存, 只是系统为了提高业务的访问IO的效率, 将读写过的文件缓存在内存中, file cache 并不会随着进程退出而释放, 只会当容器销毁或者系统内存不足时才会由系统自动回收.

但是`container_memory_working_set_bytes`的计算方式中, 只减去`total_inactive_file`的内存也是不够的, 因为还剩下`total_active_file`没有减去, 所以`container_memory_usage_bytes`反应 pod 真实的内存占用是不准确的.

## 参考文章3的操作分析

参考文章3对于自己的判断进行了认证, 作者通过`cgroup`命令手动创建资源组, 并通过`docker run`的`--cgroup-parent`选项手动指定目标资源组, 创建了一个测试容器.

对一个大文件进行`grep`, `memory.stat(total_inactive_file)`中的值就会增加相应大小的数值, 再`grep`一次, 该文件占用的内存就会转移到`memory.stat(total_active_file)`字段.

猜测是第1次`grep`进行了全文读取, kernel 便将该文件的内容全部放在了 cache, 还是`inactive`的 cache. 第2次`grep`, kernel 就将该文件的内容从`inactive`区域转移到了`active`区域, 也是比较容易理解的.

然后他通过在容器中执行`echo 3 > /proc/sys/vm/drop_caches`, 手动触发了一次内存回收, active(file) 和 inactive(file)都被回收了.
