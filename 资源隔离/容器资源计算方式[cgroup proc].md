参考文章

1. [kubernetes上报Pod已用内存不准问题分析](https://cloud.tencent.com/developer/article/1637682)
    - linux内核中, 真正意义上的"空闲内存"的计算方式, 给出了明确的公式
2. [Cgroup - Linux 内存资源管理](https://blog.csdn.net/bingqingsuimeng/article/details/52084184)
    - `free`命令中的buffer/cache:
        - buffer(Buffer Cache): 用来在系统对块设备进行读写的时候, 对块进行数据缓存的系统来使用(???没太懂)
        - cache(Page Cache): 主要用来作为文件系统上的文件数据的缓存来用, 尤其是针对当进程对文件有 read/write 操作的时候
    - cache(Page Cache)一个比较明显的特性就是, 脏数据的写回(因为内核要保证内存与设备上内容一致).
3. [值得收藏的查询进程占用内存情况方法汇总](https://www.cnblogs.com/tencentdb/articles/12272350.html)
    - golang写了几个 http 接口, 可以动态地申请内存, 然后进行分析, 挺不错的示例.
4. [记一次线上docker容器内存占用过高不释放的问题解决过程](https://www.unitimes.pro/p/801f510b5833e25886ae9bfd6b02f83e)
    - 当对文件系统进行写入时，该写入不会立即提交到磁盘，这将是非常低效的。相反，写入页面缓存的内存区域中，并且周期性地以块的形式写入磁盘
5. [docker cgroup 技术之memory（首篇）](https://www.cnblogs.com/charlieroro/p/10180827.html)
    - 从 linux 内核进程空间开始讲起(就是不太清晰明了)
6. [Docker 容器内存监控原理及应用](https://www.jb51.net/article/94677.htm)
7. [Cadvisor内存使用率指标](https://www.orchome.com/6745)
8. [Linux MemFree与MemAvailable的区别](https://blog.51cto.com/xujpxm/1961072)
    - MemTotal: 系统从加电开始到引导完成，BIOS等要保留一些内存，内核要保留一些内存，最后剩下可供系统支配的内存就是MemTotal。这个值在系统运行期间一般是固定不变的。
    - MemAvailable≈MemFree+Buffers+Cached，它是内核使用特定的算法计算出来的，是一个估计值
    - ~~OS Mem total = OS Mem used + OS Mem free~~ 我感觉这句是错的, 之后也不能肯定, 这个示例太片面

进入容器

## cgroup.procs 与 tasks

这个文件里存放着容器里所有的pid列表.

/sys/fs/cgroup/cpu/cgroup.procs

而且这2个文件在所有的 cgroup 子系统目录下都存在, 且内容都是一致的.

## /sys/fs/cgroup/memory 

### memory.stat 文件

这个文件的内容列表中, 单位都是KB

```
## page cache
cache 311427072
## 这个值与 total_rss 字段貌似完全一致, 目前没看到有什么区别.
## 表示当前pod内, **所有进程**占用的内存总量.
## 这个值基本等同于 prometheus 里 container_memory_rss 指标, 
## 也基本等同于 docker stats 的值.
## rss: Resident Set Size(常驻内存集大小), 未被放在swap空间的内存大小, 即活动内存.
rss 397410304
rss_huge 0
shmem 0
mapped_file 0
dirty 0
writeback 6488064
swap 0
pgpgin 248556
pgpgout 237150
pgfault 639936
pgmajfault 0
inactive_anon 0
active_anon 421068800
inactive_file 172097536
active_file 156762112
unevictable 0
## 貌似和 memory.limit_in_bytes 文件的值一致
hierarchical_memory_limit 3221225472
## 貌似和 memory.memsw.limit_in_bytes 文件的值一致
hierarchical_memsw_limit 3221225472
total_cache 311427072
total_rss 397410304
total_rss_huge 0
total_shmem 0
total_mapped_file 0
total_dirty 0
total_writeback 6488064
total_swap 0
total_pgpgin 248556
total_pgpgout 237150
total_pgfault 639936
total_pgmajfault 0
total_inactive_anon 0
total_active_anon 421068800
total_inactive_file 172097536
total_active_file 156762112
total_unevictable 0
```

### memory.limit_in_bytes

这个文件里存储的是当前Pod的内存Limit值. (如果不指定limit值, 貌似这个文件就会存一个大得离谱的数字...)

### memory.usage_in_bytes

prometheus 里 container_memory_usage_bytes 指标, 远大于`memory.stat(rss)`的值.

分析内核代码发现`memory.usage_in_bytes`的统计数据是包含了所有的 file cache 的, `total_active_file`和`total_inactive_file`都属于 file cache 的一部分. 并且这两个数据并不是业务真正占用的内存, 只是系统为了提高业务的访问IO的效率, 将读写过的文件缓存在内存中, file cache 并不会随着进程退出而释放, 只会当容器销毁或者系统内存不足时才会由系统自动回收.

在容器中执行`echo 3 > /proc/sys/vm/drop_caches`, 可以手动触发内存回收, `active_file`和`inactive_file`都被回收.

------

`shmget`, `shmat`等涉及的内容, 会被同时计算在 cache、mapped_file、inactive_anon 3个字段中. 比如创建100KB的共享内存, 这3个字段会同时加100KB.

与`/proc/meminfo`表示的是整个物理机的内存使用情况, 而不能反应容器本身.

- memory.limit_in_bytes: 内存使用量限制
- memory.memsw.limit_in_bytes: 内存＋ swap 空间使用的总量限制



                        ┌- memory.stat(active_file+inactive_file) file cache(可通过drop_caches回收)和buffer
                        |
                        |
memory.usage_in_bytes  -+- shmem + mlock_file: 与memory.stat(active_file+inactive_file)共同组成 cache
                        |
                        |
                        └- 

cache + buffer = memory.stat(active_file+inactive_file) + shmem + mlock_file

------

以下猜测并不准确, 但的确是在某种场景下可以对得上.

memory.stat(rss) = memory.stat(active_anon + mapped_file) (+ size of tmpfs)

memory.stat(cache) = memory.stat(active_file + inactive_file) (- size of tmpfs)
