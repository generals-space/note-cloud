# prometheus-query查询

参考文章

1. [官方文档 QUERYING PROMETHEUS](https://prometheus.io/docs/prometheus/latest/querying/basics/)

prometheus的webUI -> Graph 页面可以输入类似于sql的语句, 查询各job获得的指标信息. 

## 1. 初级过滤(瞬时值)

程序在使用prometheus提供的SDK进行日志埋点时, 会有一些通用的指标, 比如golang版本信息, 当前CPU/内存占用, GC回收情况等(这些都是可以通过profile获得), 几乎所有的`/metrics`接口都会提供这些基础信息.

那么, 当我们希望查询某个接口后的系统所使用的golang信息时, 如果不添加一些过滤条件, 将得到所有job的结果.

![](https://gitee.com/generals-space/gitimg/raw/master/94CEF94EB2114279957792C443A3DDD6.png)

> 注意: 不同程序在使用`metric`进行埋点时都会指定一个名称前缀, 比如`prometheus`用于监控自身的指标前缀都为`prometheus_`, kubernetes的则一般为`kube_`. 当你不熟悉 prometheus 的界面时, 你会发现 Graph 页面中, `Execute`按钮后面的下拉框里有超~~多的指标, 都不知道自己要找哪个.

我们可以先根据这全部的结果, 以及结果中大括号`{}`中参数进行进一步的过滤.

上面的结果中只有一行, 是因为当前prometheus中只有一个job, 抓取的是prometheus本身的指标信息. 如果是多个job, 那么可以尝试将过滤条件写作

![](https://gitee.com/generals-space/gitimg/raw/master/16E6CCDBD225D03C4A78523B313550CF.png)

> `go_info`与`go_info{}`的作用一样, 查询全部结果.

注意: 查询结果中, 右侧的`Value`列即为本次查询的**瞬时值**.

对于某一指标各字段的过滤方法, 不只可以用`=`, 还有其他如`!=`, `=~`等, 详细介绍可以见参考文章1.

`Graph/Console`标签页下面都有可以选择的时间点, 以查询历史时间中的瞬时值, 但是没有范围变动, 所以在`Graph`界面中你看到的大多是一条/n条横线而已, 没有起伏变化.

我之后又尝试查看了其他的指标, 发现结果表格中"Element"列可供使用的过滤字段只有`instance`和`job`两个, 猜测应该是因为prometheus这个job只有这两个label吧.

![](https://gitee.com/generals-space/gitimg/raw/master/B4B4C41C791DAC5F9C5E88F104DD5ED7.png)

## 2. 时间过滤(范围值)

关于时间过滤有两个参数, 一个是范围, 一个是偏移量 offset.

### 