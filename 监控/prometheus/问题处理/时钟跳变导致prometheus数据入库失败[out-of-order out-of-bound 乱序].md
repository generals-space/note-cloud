参考文章

1. [Prometheus raise out of bounds error for all targets after resume the linux system from a suspend #8243](https://github.com/prometheus/prometheus/issues/8243)
    - 问题日志与我遇到的相似.
2. [fix prometheus out of bound bug #12182](https://github.com/prometheus/prometheus/issues/12182)
    - python写的修复脚本
3. [out of bound fix #12183](https://github.com/prometheus/prometheus/issues/12183)
    - 同参考文章2
    - 官方提供了一个"out-of-order sample"特性, 可以在不删除数据的情况下恢复.
4. [Add out-of-order sample support to the TSDB #11075](https://github.com/prometheus/prometheus/pull/11075)
    - 参考文章3提供的特性链接, 2.39.x版本提供了"storage.tsdb.out_of_order_time_window"参数.
5. [[tsdb] Ingest out of order samples and samples from a few hours ago #8535](https://github.com/prometheus/prometheus/issues/8535)
6. [prometheus configuration tsdb](https://prometheus.io/docs/prometheus/2.53/configuration/configuration/#tsdb)
    - 官方文档配置
7. [如何解决Prometheus的数据回填问题](https://blog.csdn.net/sinat_32582203/article/details/128727107)
8. [记一次远程写性能问题引发的Prometheus版本升级事件](https://cloud.tencent.com/developer/article/2314673)
    - `storage.tsdb.out_of_order_time_window`在配置文件中写法, `storage`与`global`平级

## 场景描述

prometheus(镜像版本): v2.53.2

生产环境中可能因为 ntp server 同步异常导致系统时钟跳变, 测试同学测了个极端场景, 手动使用`date -s`将系统时间向后调1年, 再恢复, 然后发现prometheus不再采集数据了.

查看 prometheus 日志有如下输出

```log
Dec 01 09:46:47 prometheus[1629652]: level=warn ts=2020-12-01T01:46:47.293Z caller=scrape.go:1378 component="scrape manager" scrape_pool=ssh target="http://127.0.0.1:9115/probe?module=ssh_banner&target=172.20.149.141%3A22" msg="Error on ingesting samples that are too old or are too far into the future" num_dropped=6
Dec 01 09:46:47 prometheus[1629652]: level=warn ts=2020-12-01T01:46:47.293Z caller=scrape.go:1145 component="scrape manager" scrape_pool=ssh target="http://127.0.0.1:9115/probe?module=ssh_banner&target=xxxx%3A22" msg="Append failed" err="out of bounds"
Dec 01 09:46:47 prometheus[1629652]: level=warn ts=2020-12-01T01:46:47.293Z caller=scrape.go:1094 component="scrape manager" scrape_pool=ssh target="http://127.0.0.1:9115/probe?module=ssh_banner&target=xxx%3A22" msg="Appending scrape report failed" err="out of bounds"
```

其实在系统时间向后调1年的时候, prometheus就会因为"too old or are too far into the future"无法正常工作了, 入库时发现当前时间与上一条记录的时间相差过大而丢弃.

测试要求时间恢复后所有组件仍可正常运行.

但是在时间恢复后仍然会报"out of bounds(越界)", 无法写入数据, 只能清空数据目录, 然后重启prometheus服务才可以.

## 问题分析

采集端prometheus时间向后调1小时(或更久)就会出现如下报错, 不再采集数据了.

```log
ts=2025-04-11T07:35:26.073Z caller=scrape.go:1747 level=warn component="scrape manager" scrape_pool=k8s_cAdvisor target=https://192.168.203.253:6443/api/v1/nodes/white-master-2/proxy/metrics/cadvisor msg="Error on ingesting samples that are too old or are too far into the future" num_dropped=19128
ts=2025-04-11T07:35:33.005Z caller=compact.go:576 level=info component=tsdb msg="write block" mint=1744344000009 maxt=1744351200000 ulid=01JRHWCY05B5CT0JZ4YVV4HYMV duration=1.416416456s ooo=false
ts=2025-04-11T07:35:33.163Z caller=head.go:1355 level=info component=tsdb msg="Head GC completed" caller=truncateMemory duration=155.613314ms
ts=2025-04-11T07:35:35.002Z caller=scrape.go:1747 level=warn component="scrape manager" scrape_pool=k8s_cAdvisor target=https://192.168.203.253:6443/api/v1/nodes/white-master-3/proxy/metrics/cadvisor msg="Error on ingesting samples that are too old or are too far into the future" num_dropped=22670
ts=2025-04-11T07:35:39.562Z caller=scrape.go:1747 level=warn component="scrape manager" scrape_pool=k8s_cAdvisor target=https://192.168.203.253:6443/api/v1/nodes/white-master-1/proxy/metrics/cadvisor msg="Error on ingesting samples that are too old or are too far into the future" num_dropped=19502
```

但其实调整10分钟往上, 虽然不报错, 但采集也会出问题.

假设当前时间为13:00, 调整的时间为10分钟后, 13:10, 当真实世界的时间来到13:10后, 13:00-13:10这10分钟间的数据为空白, 而13:10-13:20这10分钟的数据会因发生重叠而被丢弃.

```
    13:00    13:10
     ↓        ↓
======        ======
```

数据重叠

```log
ts=2025-04-12T05:50:12.364Z caller=scrape.go:1744 level=warn component="scrape manager" scrape_pool=controller_metrics target=http://192.168.11.2:6798/metrics msg="Error on ingesting samples with different value but same timestamp" num_dropped=332
```

如果在10分钟内将时间调整为正常值, 则有采集指标会出现乱序.

```
    13:00    13:10
     ↓        ↓
======        ======
          ===========
          ↑
          13:05
```

乱序

```log
ts=2025-04-12T05:31:56.394Z caller=scrape.go:1741 level=warn component="scrape manager" scrape_pool=k8s_cAdvisor target=https://192.168.203.253:6443/api/v1/nodes/white-master-2/proxy/metrics/cadvisor msg="Error on ingesting out-of-order samples" num_dropped=434
```

乱序也会导致新数据无法入库, 但由于旧数据

## 解决方案

```yaml
global:
...
storage:
  tsdb:
     out_of_order_time_window: 12h
```

修改配置前的日志

```log
ts=2025-04-10T05:40:28.311Z caller=scrape.go:1747 level=warn component="scrape manager" scrape_pool=prometheus target=http://localhost:9090/metrics ms
g="Error on ingesting samples that are too old or are too far into the future" num_dropped=552
ts=2025-04-10T05:40:28.311Z caller=scrape.go:1294 level=warn component="scrape manager" scrape_pool=prometheus target=http://localhost:9090/metrics ms
g="Appending scrape report failed" err="out of bounds"
ts=2025-04-10T05:40:30.814Z caller=scrape.go:1747 level=warn component="scrape manager" scrape_pool=controller_metrics target=http://192.168.11.2:6798
/metrics msg="Error on ingesting samples that are too old or are too far into the future" num_dropped=16
ts=2025-04-10T05:40:30.814Z caller=scrape.go:1294 level=warn component="scrape manager" scrape_pool=controller_metrics target=http://192.168.11.2:6798
/metrics msg="Appending scrape report failed" err="out of bounds"
```

添加配置后重启

```log
ts=2025-04-10T05:40:42.785Z caller=web.go:568 level=info component=web msg="Start listening for connections" address=0.0.0.0:9090
ts=2025-04-10T05:40:42.785Z caller=main.go:1148 level=info msg="Starting TSDB ..."
ts=2025-04-10T05:40:42.787Z caller=repair.go:56 level=info component=tsdb msg="Found healthy block" mint=1744250764649 maxt=1744257600000 ulid=01JRF1ENHAHN7S0DM5KDY4Q1TC
ts=2025-04-10T05:40:42.787Z caller=repair.go:56 level=info component=tsdb msg="Found healthy block" mint=1744257600000 maxt=1744264800000 ulid=01KNTX37K1K2RFGCVE9CWDTV18

4-10T05:40:42.821Z caller=main.go:1133 level=info msg="Server is ready to receive web requests."
ts=2025-04-10T05:40:42.821Z caller=manager.go:164 level=info component="rule manager" msg="Starting rule manager..."
ts=2025-04-10T05:40:58.363Z caller=compact.go:567 level=info component=tsdb msg="write block resulted in empty block" mint=1744264800000 maxt=1744272000000 duration=31.893535ms
ts=2025-04-10T05:40:58.374Z caller=head.go:1355 level=info component=tsdb msg="Head GC completed" caller=truncateMemory duration=7.27357ms
ts=2025-04-10T05:40:58.411Z caller=compact.go:567 level=info component=tsdb msg="write block resulted in empty block" mint=1744272000000 maxt=1744279200000 duration=37.449207ms

...省略
ts=2025-04-10T05:43:09.151Z caller=head.go:1355 level=info component=tsdb msg="Head GC completed" caller=truncateMemory duration=8.321722ms
ts=2025-04-10T05:43:09.152Z caller=checkpoint.go:101 level=info component=tsdb msg="Creating checkpoint" from_segment=5 to_segment=6 mint=1775793600000
ts=2025-04-10T05:43:09.167Z caller=head.go:1317 level=info component=tsdb msg="WAL checkpoint complete" first=5 last=6 duration=15.568347ms
ts=2025-04-10T05:43:09.223Z caller=compact.go:576 level=info component=tsdb msg="write block" mint=1744257600000 maxt=1744264800000 ulid=01JRF3JEKREMXTW6ST878Y6G0G duration=47.182091ms ooo=true

ts=2025-04-10T05:43:09.223Z caller=db.go:1356 level=info component=tsdb msg="out-of-order compaction completed" duration=47.323207ms ulids=[01JRF3JEKREMXTW6ST878Y6G0G]
ts=2025-04-10T05:43:09.224Z caller=db.go:1549 level=warn component=tsdb msg="Overlapping blocks found during reloadBlocks" detail="[mint: 1744257600000, maxt: 1744264800000, range: 2h0m0s, blocks: 2]: <ulid: 01JRF3JEKREMXTW6ST878Y6G0G, mint: 1744257600000, maxt: 1744264800000, range: 2h0m0s>, <ulid: 01KNTX37K1K2RFGCVE9CWDTV18, mint: 1744257600000, maxt: 1744264800000, range: 2h0m0s>"
ts=2025-04-10T05:43:09.232Z caller=head.go:1355 level=info component=tsdb msg="Head GC completed" caller=truncateOOO duration=7.328265ms
ts=2025-04-10T05:43:09.248Z caller=compact.go:762 level=info component=tsdb msg="Found overlapping blocks during compaction" ulid=01JRF3JENGN43N9E1M8EECX1N6
ts=2025-04-10T05:43:09.300Z caller=compact.go:514 level=info component=tsdb msg="compact blocks" count=2 mint=1744257600000 maxt=1744264800000 ulid=01JRF3JENGN43N9E1M8EECX1N6 sources="[01JRF3JEKREMXTW6ST878Y6G0G 01KNTX37K1K2RFGCVE9CWDTV18]" duration=68.180106ms
ts=2025-04-10T05:43:09.303Z caller=db.go:1712 level=info component=tsdb msg="Deleting obsolete(过时的) block" block=01KNTX37K1K2RFGCVE9CWDTV18
ts=2025-04-10T05:43:09.304Z caller=db.go:1712 level=info component=tsdb msg="Deleting obsolete(过时的) block" block=01JRF3JEKREMXTW6ST878Y6G0G
```

~~删除了过时的块, 还是会把旧数据清理掉, 这个方法没意义啊...~~ 哦, 是因为prometheus的`--storage.tsdb.retention.tim`参数只设置了7天, 即只保留7天的数据, 调整时间到1年后就会把旧数据清理掉. 如果只向后调1小时, 就不会全部清理了.

------

添加此配置后, 调整时间到1小时后, 再修改回来, 还是不能正常采集数据, 需要手动重启prometheus或是reload才可以.

猜测: 开启此配置并重启后, 就相当于以当前时间为真实时间, 在允许的时间范围内直接进行覆写, 也不管是否冲突了.
