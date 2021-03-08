# Docker的资源隔离

参考文章

1. [Docker 运行时资源限制](https://blog.csdn.net/candcplusplus/article/details/53728507)
    - docker对内存限制的几个特点
    - linux的OOM机制
    - 核心内存与用户内存的概念
    - CPU限制的可用选项
    - CPU资源的相对限制与绝对限制
2. [Docker容器资源限制测试](https://www.centos.bz/2017/12/docker%E5%AE%B9%E5%99%A8%E8%B5%84%E6%BA%90%E9%99%90%E5%88%B6%E6%B5%8B%E8%AF%95/)
    - stress压测工具测试docker的资源限制效果
3. [[经验分享] docker的资源隔离---cpu、内存、磁盘限制](http://www.iyunv.com/thread-116572-1-1.html)
4. [官方文档 Managing Compute Resources for Containers](https://kubernetes.io/docs/concepts/configuration/manage-compute-resources-container/)
    - 三种可被管理的资源: CPU, Memory, HugePage

`spec.containers[].resources.requests.cpu`有两种格式: `0.5`与`500m`相同, 不过由于小数精度的问题, 建议使用后者.

`500m`这种格式中的`m`表示`millicpu/millicores`, 可以理解为CPU算力, 表示可以保证单核主机上1/2的运算能力. 且该值是绝对值, 即总量而不是百分比. 在双核主机上, `500m`就只能占用1/4的运算能力了.

```
      --cpu-period int                 Limit CPU CFS (Completely Fair Scheduler) period
      --cpu-quota int                  Limit CPU CFS (Completely Fair Scheduler) quota
      --cpu-rt-period int              Limit CPU real-time period in microseconds
      --cpu-rt-runtime int             Limit CPU real-time runtime in microseconds
  -c, --cpu-shares int                 CPU shares (relative weight)
      --cpus decimal                   Number of CPUs
      --cpuset-cpus string             CPUs in which to allow execution (0-3, 0,1)
      --cpuset-mems string             MEMs in which to allow execution (0-3, 0,1)
```

```
  -m, --memory bytes                   Memory limit
      --memory-reservation bytes       Memory soft limit
      --memory-swap bytes              Swap limit equal to memory plus swap: '-1' to enable unlimited swap
      --memory-swappiness int          Tune container memory swappiness (0 to 100) (default -1)
      --mount mount                    Attach a filesystem mount to the container
```
