参考文章

1. [prometheus-book](https://yunlzheng.gitbook.io/prometheus-book/)
    - 一个完善的监控目标是要能够从白盒的角度发现潜在问题, 能够在黑盒的角度快速发现已经发生的问题(第4章 [网络探测：Blackbox Exporter]).
2. [container-monitor](https://yasongxu.gitbook.io/container-monitor/)
3. [监控神器Prometheus用不对，也就是把新手村的剑](https://cloud.tencent.com/developer/article/1660745)

metrics: 指标

prometheus 通过 job 配置指定要监测的对象, 每次更新 job 配置都需要重启 prometheus.

