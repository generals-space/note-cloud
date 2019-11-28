每个node可能有3种角色: 

1. control plane
2. etcd
3. worker

单节点时这3个都要.

rke要求部署节点时使用非root用户的权限, 否则ssh将失败.

> FATA[0000] Cluster must have at least one etcd plane host: failed to connect to the following etcd host(s) [192.168.0.211]

```
user add ubuntu
usermod -aG docker ubuntu
```

