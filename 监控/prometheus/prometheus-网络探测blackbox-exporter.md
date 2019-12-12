# prometheus-网络探测blackbox-exporter

参考文章

1. [网络探测：Blackbox Exporter](https://yunlzheng.gitbook.io/prometheus-book/part-ii-prometheus-jin-jie/exporter/commonly-eporter-usage/install_blackbox_exporter)
    - 介绍了blackbox_exporter可执行文件的配置及使用方法.
2. [对Ingress和Service进行网络探测](https://yunlzheng.gitbook.io/prometheus-book/part-iii-prometheus-shi-zhan/readmd/use-prometheus-monitor-kubernetes#dui-ingress-he-service-jin-hang-wang-luo-tan-ce)
    - 介绍了在kubernetes中通过`blackbox_exporter`自动发现service和ingress资源及监控的方法.

我们知道, prometheus可以抓取指定目标的`/metrics`接口, 以获取大量运行中的状态信息. 但这些接口都是需要在程序中添加埋点才能处理的, 但如果是非golang, 非通用的系统, 没有提供`/metrics`(比如nginx, redis)怎么办? 在prometheus中如何获取这种类型的健康状态?

此时就需要`blackbox-exporter`组件了, ta提供了传统的, 基于IP和端口的监控方式. ta提供了icmp模块检测主机存活, tcp模块检测端口开放, http模块检测restful服务等.

按照参考文章1和2, 在集群中部署此服务将非常容易. 

在prometheus本身提供的监控服务中, `/metrics`接口的确只能返回2xx, 400/404这些都会将目标视为`DOWN`的状态. blackbox-exporter的默认配置也是如此. 如下

```yaml
modules:
  http_2xx:
    prober: http
  http_post_2xx:
    prober: http
    http:
      method: POST
```

但我们可以对其稍作修改, 接受3xx,4xx的状态码, 也可以调整响应的超时时间, 定义请求头等. 可以见官方文档的示例配置, 非常易懂, 可以说判活的操作十分全面了.
