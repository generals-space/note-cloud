
## 场景一: besteffort级的cgroup设置的cpu 内存使用上限是否可限制住下级的任务进程

结论: 可以. 
cpu设置值将由besteffort下所有容器争抢, 内存则由所有容器共用(超出设置上限会引发OOM)

## 场景二: 若单个子任务资源使用超过新设置的整体资源上限，是否可成功设置

结论: CPU可以, 内存不可以

CPU可以设置为小于当前besteffort级所有容器所使用的核数总量, 之后各容器将在该值限制下进行争抢, 占用较多的进程CPU使用率会下降到设置值下的合理区间;
当设置的内存值小于当前besteffort级所有容器所使用的内存总量时, 会失败, 报"Device or resource busy"

## besteffort下单个子cgroup资源设置上限是否可超过besteffort级cgroup资源上限

结论: cpu不可以, 内存可以但无效.

CPU在子级目录设置超过父级目录的设置值时, 会出现错误: "write error: Invalid argument"
