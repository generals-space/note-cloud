# prometheus热重启[reload]

参考文章

1. [Configuration](https://prometheus.io/docs/prometheus/2.53/configuration/configuration/)

需要prometheus在启动时添加了`--web.enable-lifecycle`选项.

```
curl -XPOST localhost:9090/-/reload
```
