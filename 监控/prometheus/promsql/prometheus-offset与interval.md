参考文章

1. [Prometheus学习（六）之Prometheus查询说明](https://www.cnblogs.com/even160941/p/15109453.html)
2. [Prometheus监控系统 安装与配置详细教程](https://www.cnblogs.com/ExMan/p/12567247.html)

offset 就像 grafana 右上角选择的, [5分钟内, 10分钟内, 1小时内, 1天内, 1周内, 一个月内]的阶段框选.

比如, 如下语句查看的是, 一天内的get请求总量.

```
sum(http_requests_total{method="GET"} offset 1d)
```

interval 则为数据汇总的分段值.

在查询语句中, 写作如下

```
http_requests_total{job="prometheus"}[5m]
```
