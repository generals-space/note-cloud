# prometheus的四种指标类型[counter gauge histogram summary]

参考文章

1. [Prometheus 系统监控方案 一](https://www.cnblogs.com/vovlie/p/Prometheus_CONCEPTS.html)
    - 简要介绍了下 prometheus 的4种指标类型, 对 counter, gauge 的解释不错, histogram, summary 的解释不太行
2. [理解时间序列](https://www.cnblogs.com/sanduzxcvbnm/p/13306085.html)
3. [Prometheus 指标类型](https://fuckcloudnative.io/prometheus/2-concepts/metric_types.html)
    - 由API请求的"长尾问题", 得出`histogram`直方图的必要性
    - Prometheus 中的 histogram 是累积直方图
    - quantile 分位数
4. [一文搞懂 Prometheus 的直方图](https://www.jianshu.com/p/ac7ab610537a/)
    - 值得一看
5. [官方文档 HISTOGRAMS AND SUMMARIES](https://prometheus.io/docs/practices/histograms/)
6. [如何通俗地理解分位数？](https://www.zhihu.com/question/67763556)
    - 中位数(1/2分位数)
7. [分布统计：Heatmap面板](https://yunlzheng.gitbook.io/prometheus-book/part-ii-prometheus-jin-jie/grafana/grafana-panels/use_heatmap_panel)
    - Graph面板来可视化Histogram类型的监控指标不够直观, 最好使用Heatmap面板

一般指标项大多为`gauge`, 表示某个瞬时值, 如当前工程中的线程数, 正在处理的请求数等, 是不断变动的. 由prometheus按时间序列采集后, 在grafana里直接写上这类指标, 就能得到由点组成的折线.

counter则为一个累加值, 如自从工程启动, 经历的GC次数, 处理过的请求总数等. counter值可以根据rate()函数计算出累似于gauge的"速率"数据.

而 histogram 和 summary 就不太容易理解了, 建议查看参考文章4.

## Prometheus 中的 histogram 是累积直方图

```
# HELP controller_runtime_reconcile_time_seconds Length of time per reconciliation per controller
# TYPE controller_runtime_reconcile_time_seconds histogram
controller_runtime_reconcile_time_seconds_bucket{controller="kafkacluster",le="0.005"} 16211
controller_runtime_reconcile_time_seconds_bucket{controller="kafkacluster",le="0.01"} 28558
controller_runtime_reconcile_time_seconds_bucket{controller="kafkacluster",le="0.025"} 36694
controller_runtime_reconcile_time_seconds_bucket{controller="kafkacluster",le="0.05"} 36821
controller_runtime_reconcile_time_seconds_bucket{controller="kafkacluster",le="0.1"} 36946
controller_runtime_reconcile_time_seconds_bucket{controller="kafkacluster",le="0.25"} 37062
controller_runtime_reconcile_time_seconds_bucket{controller="kafkacluster",le="0.5"} 37137
controller_runtime_reconcile_time_seconds_bucket{controller="kafkacluster",le="1"} 37150
controller_runtime_reconcile_time_seconds_bucket{controller="kafkacluster",le="+Inf"} 37150
controller_runtime_reconcile_time_seconds_sum{controller="kafkacluster"} 273.2828634999991
controller_runtime_reconcile_time_seconds_count{controller="kafkacluster"} 37150
```

kubebuilder工程中, `reconciler()`函数的处理耗时统计, 从程序启动, 已经进行了37150次reconcile操作(此值应该类似于`counter`, 是一个累加值, 单调递增, 不会减少), 这所有的reconcile操作共耗时`273.2828634999991`秒.

在这些reconcile操作中, 有16211个操作耗时小于0.005秒, 有28558个操作耗时小于0.01秒(注意这里包含了前面的16211个操作), 有36694个操作耗时小于0.025秒(包含了前面的28558个操作), 换成表格表示如下

| 耗时       | 处理数量 |
| :--------- | :------- |
| (0, 0.005) | 16211    |
| (0, 0.01)  | 28558    |
| (0, 0.025) | 36694    |
| (0, 0.05)  | 36821    |
| (0, 0.1)   | 36946    |
| (0, 0.25)  | 37062    |
| (0, 0.5)   | 37137    |
| (0, 0.1)   | 37150    |
| (0, ∞)     | 37150    |

这种每次都从0开始计数, 而不是像(0, 0.005), (0.005, 0.01), (0.01, 0.025)这种分段式区间的表示形式, 被称为"累积直方图".

## quantile 分位数

参考文章3介绍了一个新的概念: **分位数**. 

按照参考文章6的说法, 中位数就是1/2的分位数, 其他的像1/4分位数会有一前一后两个.

```
# HELP go_gc_duration_seconds A summary of the GC invocation durations.
# TYPE go_gc_duration_seconds summary
go_gc_duration_seconds{quantile="0"} 0
go_gc_duration_seconds{quantile="0.25"} 0
go_gc_duration_seconds{quantile="0.5"} 0
go_gc_duration_seconds{quantile="0.75"} 0
go_gc_duration_seconds{quantile="1"} 0.0024302
go_gc_duration_seconds_sum 0.2147324
go_gc_duration_seconds_count 1762
```

上面的`quantile`, 以及上一节的`le(less than)`都可以认为是实际的分位数值.

## 直方图的panel展现形式

参考文章7中提到, 使用`Graph`面板虽然可以展示Histogram的指标, 但是不够直观, 如下

![](https://gitee.com/generals-space/gitimg/raw/master/eb3cd2d627e862de4c608facfba81e45.png)

![](https://gitee.com/generals-space/gitimg/raw/master/7b4e1bec6f891a3702e6a0e1d78cc625.png)

只能看到各项区间(各区间的折线并不是累积的, 而是求的差值, 挺贴心的)的增长速度, 看不到分布情况.

