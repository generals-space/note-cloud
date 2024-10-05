参考文章

1. [容器化Node Exporter对主机磁盘使用率的监控问题](https://www.cnblogs.com/YaoDD/p/12357169.html)

```yaml
    volumeMounts:
    - mountPath: /var/lib/cni/sriov
      mountPropagation: HostToContainer
      name: host-var-lib-cni-sriov
```
