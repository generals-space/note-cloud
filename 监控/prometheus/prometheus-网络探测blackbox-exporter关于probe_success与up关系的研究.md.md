# prometheus-网络探测blackbox-exporter关于probe_success与up关系的研究.md

参考文章

1. [Difference between probe_success and up?](https://stackoverflow.com/questions/51984837/difference-between-probe-success-and-up)
2. [Prometheus官方文档 Failed scrapes](https://prometheus.io/docs/instrumenting/writing_exporters/#failed-scrapes)
3. [Checking for HTTP 200s with the Blackbox Exporter](https://www.robustperception.io/checking-for-http-200s-with-the-blackbox-exporter)

```
Metrics that would have been returned:
...省略
probe_http_content_length 0
...省略
probe_http_status_code 0
...省略
probe_success 0
```
