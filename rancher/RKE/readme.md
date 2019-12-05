参考文章

1. [General Linux Requirements](https://rancher.com/docs/rke/latest/en/os/#general-linux-requirements)

rke版本: 1.0.0

当前目录的`cluster.yml`为`rke up`所需要的配置文件, 在各节点安装好依赖后可根据此文件启动安装流程. `cluster_full.yml`则是使用`rke config`通过交互式命令生成的配置, 交互的内容可见`rke_config.txt`. `cluster_full.yml`只填写了必填字段, 还有很多是留空的, 不适合备份与重用.

每个node可能有3种角色: 

1. control plane
2. etcd
3. worker

单节点时这3个都要.

rke要求部署节点时使用非root用户的权限, 否则ssh将失败. 见参考文章1.

```
WARN[0000] Failed to set up SSH tunneling for host [192.168.0.211]: Can't retrieve Docker Info: Cannot connect to the Docker daemon at unix:///var/run/docker.sock. Is the docker daemon running?
...省略
FATA[0000] Cluster must have at least one etcd plane host: failed to connect to the following etcd host(s) [192.168.0.211]
```

```
user add ubuntu
usermod -aG docker ubuntu
```

`rke up`之后会生成2个额外文件: 

1. `kube_config_cluster`: `kubectl`配置文件.
2. `cluster.rkestate`: 包含整个集群所有的配置, 包括所有组件的证书和密钥, 比前者内容要丰富得多.

