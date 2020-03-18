参考文章

1. [Kubernetes 通过statefulset部署redis cluster集群](https://www.cnblogs.com/kuku0223/p/10906003.html)
    - statefulset+NFS存储部署redis-cluster集群示例

rancher的local-path不支持`readWriteMany`访问模式, 所以这里我们使用`nfs-provisioner`存储.

在statefulset创建的pod资源全部启动后, 需要使用`redis-trib`进行部署, 参考文章中的示例命令使用`dig`查找各redis实例的IP, 效果如下

```console
$ dig +short redis-app-0.redis-service.default.svc.cluster.local
10.254.0.159
```

> `redis-trib`不支持使用域名来创建集群.

参考文章1中的示例命令如下.

```bash
redis-trib create --replicas 1 \
`dig +short redis-app-0.redis-service.default.svc.cluster.local`:6379 \
`dig +short redis-app-1.redis-service.default.svc.cluster.local`:6379 \
`dig +short redis-app-2.redis-service.default.svc.cluster.local`:6379 \
`dig +short redis-app-3.redis-service.default.svc.cluster.local`:6379 \
`dig +short redis-app-4.redis-service.default.svc.cluster.local`:6379 \
`dig +short redis-app-5.redis-service.default.svc.cluster.local`:6379
```

个人觉得太长了, 想简化一点, 步骤如下

```console
$ dig +short redis-app-{0..5}.redis-service.default.svc.cluster.local
10.254.0.159
10.254.0.161
10.254.0.163
10.254.0.165
10.254.0.167
10.254.0.169
```

```console
$ dig +short redis-app-{0..5}.redis-service.default.svc.cluster.local | awk '{printf("%s:6379 ", $1)}'
10.254.0.159:6379 10.254.0.161:6379 10.254.0.163:6379 10.254.0.165:6379 10.254.0.167:6379 10.254.0.169:6379
```

> 这里用到了`awk`的格式化输出.

最终使用的命令如下

```console
$ clusterIPs=$(dig +short redis-app-{0..5}.redis-service.default.svc.cluster.local | awk '{printf("%s:6379 ", $1)}')
$ redis-trib create --replicas 1 $clusterIPs
>>> Creating cluster
>>> Performing hash slots allocation on 6 nodes...
Using 3 masters:
10.254.0.159:6379
10.254.0.161:6379
10.254.0.163:6379
Adding replica 10.254.0.165:6379 to 10.254.0.159:6379
Adding replica 10.254.0.167:6379 to 10.254.0.161:6379
Adding replica 10.254.0.169:6379 to 10.254.0.163:6379
M: 9ce9f7165390291cf62d960c0392b8fbc4afc5ed 10.254.0.159:6379
   slots:0-5460 (5461 slots) master
M: 2c1bc368011482940e1561cb198aa00f2c1fd599 10.254.0.161:6379
   slots:5461-10922 (5462 slots) master
M: 03d9c220724d424c63a8476bab3fc6cee2230462 10.254.0.163:6379
   slots:10923-16383 (5461 slots) master
S: 72b933cb897f1a464bce52a3335ae7be7a2db4cc 10.254.0.165:6379
   replicates 9ce9f7165390291cf62d960c0392b8fbc4afc5ed
S: fba2e02de8bac1181e80f71e470b388c432ee004 10.254.0.167:6379
   replicates 2c1bc368011482940e1561cb198aa00f2c1fd599
S: 805d2680e2623681e92b80f729dcae50340e44cf 10.254.0.169:6379
   replicates 03d9c220724d424c63a8476bab3fc6cee2230462
Can I set the above configuration? (type 'yes' to accept): yes                ## 这里需要用户自行确认.
>>> Nodes configuration updated
>>> Assign a different config epoch to each node
>>> Sending CLUSTER MEET messages to join the cluster
Waiting for the cluster to join.....
>>> Performing Cluster Check (using node 10.254.0.159:6379)
M: 9ce9f7165390291cf62d960c0392b8fbc4afc5ed 10.254.0.159:6379
   slots:0-5460 (5461 slots) master
   1 additional replica(s)
M: 2c1bc368011482940e1561cb198aa00f2c1fd599 10.254.0.161:6379@16379
   slots:5461-10922 (5462 slots) master
   1 additional replica(s)
S: 805d2680e2623681e92b80f729dcae50340e44cf 10.254.0.169:6379@16379
   slots: (0 slots) slave
   replicates 03d9c220724d424c63a8476bab3fc6cee2230462
S: fba2e02de8bac1181e80f71e470b388c432ee004 10.254.0.167:6379@16379
   slots: (0 slots) slave
   replicates 2c1bc368011482940e1561cb198aa00f2c1fd599
S: 72b933cb897f1a464bce52a3335ae7be7a2db4cc 10.254.0.165:6379@16379
   slots: (0 slots) slave
   replicates 9ce9f7165390291cf62d960c0392b8fbc4afc5ed
M: 03d9c220724d424c63a8476bab3fc6cee2230462 10.254.0.163:6379@16379
   slots:10923-16383 (5461 slots) master
   1 additional replica(s)
[OK] All nodes agree about slots configuration.
>>> Check for open slots...
>>> Check slots coverage...
[OK] All 16384 slots covered.
```

> `--replicas 1`: 创建的集群中为每个主节点分配一个从节点, 上面的6个节点可以达到3主3从.

接下来是验证阶段.

1. 进入任意redis实例命令行中, 使用`cluster info`命令可以查看集群信息.
2. 到挂载目录去查看是否生成了数据文件.

------

接下来验证cluster的高可用和动态扩容特性.

## 主从切换

## 动态扩容
