# prometheus无法采集controller和scheduler的状态信息

参考文章

1. [全手动部署prometheus-operator监控Kubernetes集群遇到的坑](https://www.servicemesher.com/blog/prometheus-operator-manual/)
    - 默认kuber中controller和scheduler没有svc, prometheus无法采集数据的坑

```yaml
      kubernetes_sd_configs:
      - role: endpoints
        ## 注意: apiserver的ep是在default命名空间下的, 而不是kube-system
        namespaces:
          names:
          - default
```

在使用kuber的服务发现时, 如果role不指定命名空间, 那么搜索范围会是全部ns, 可能会在prometheus的dashboard中看到许多不相关的记录(而且还是down状态)来干扰.
