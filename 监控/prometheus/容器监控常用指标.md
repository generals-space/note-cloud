# 容器监控常用指标

参考文章

1. [prometheus 常用指标](https://www.jianshu.com/p/5602af08f432)
2. [cadvisor metrics container_memory_working_set_bytes vs container_memory_usage_bytes](https://blog.csdn.net/u010918487/article/details/106190764)
    - [How much is too much? The Linux OOMKiller and “used” memory](https://medium.com/faun/how-much-is-too-much-the-linux-oomkiller-and-used-memory-d32186f29c9d)
    - [A Deep Dive into Kubernetes Metrics — Part 3 Container Resource Metrics](https://blog.freshtracks.io/a-deep-dive-into-kubernetes-metrics-part-3-container-resource-metrics-361c5ee46e66)

## 磁盘IO

```
sum(rate(container_fs_writes_bytes_total[5m])) by (container_name,device)
sum(rate(container_fs_reads_bytes_total[5m])) by (container_name,device)
```

## 网络IO

```
sum(rate(container_network_receive_bytes_total[5m])) by (name)
sum(rate(container_network_transmit_bytes_total[5m])) by (name)
```
