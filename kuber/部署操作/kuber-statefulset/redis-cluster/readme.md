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
clusterIPs=$(dig +short redis-app-{0..5}.redis-service.default.svc.cluster.local | awk '{printf("%s:6379 ", $1)}')
redis-trib create --replicas 1 $clusterIPs
```

> `--replicas 1`: 创建的集群中为每个主节点分配一个从节点, 上面的6个节点可以达到3主3从.

打印结果如下

```
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

接下来是验证阶段.

1. 进入任意redis实例命令行中, 使用`cluster info`命令可以查看集群信息.
2. 到挂载目录去查看是否生成了数据文件.

------

接下来验证cluster的高可用和动态扩容特性.

## 主从切换

以`redis-app-1`实例为例, 集群创建完成后其角色为`master`.

```
127.0.0.1:6379> role
1) "master"
2) (integer) 19418
3) 1) 1) "10.254.1.11"
      2) "6379"
      3) "19418"
```

然后删除ta.

> 虽然说是`statefulset`在创建时会按顺序, 但是我在删除`redis-app-1`时, 也只是重建了ta本身而已, 并没有影响后面的2,3,4,5.

然后其角色变为了`slave`

```
127.0.0.1:6379> role
1) "slave"
2) "10.254.1.11"
3) (integer) 6379
4) "connected"
5) (integer) 20090
```

## 动态扩容

由于我们使用的是NFS的provisioner, 所以无需再手动创建NFS目录和pvc资源了, 可以直接扩容.

将`02.ss.yaml`中的`replicas`字段修改为8, 重新`apply`.

这会新增两个节点, 但是这两个节点并没有加入到集群中(在0-5实例中使用`cluster info`查看集群状态, 会发现`cluster_known_nodes`还是6), 需要再次使用`redis-trib`工具将新节点加入.

```console
$ redis-trib add-node \
$(dig +short redis-app-6.redis-service.default.svc.cluster.local):6379 \
$(dig +short redis-app-0.redis-service.default.svc.cluster.local):6379

$ redis-trib add-node \
$(dig +short redis-app-7.redis-service.default.svc.cluster.local):6379 \
$(dig +short redis-app-0.redis-service.default.svc.cluster.local):6379
```

`add-node`后面跟的是新节点的信息, 再后面是以前集群中的任意一个节点.

打印结果如下

```
>>> Adding node 10.254.0.194:6379 to cluster 10.254.0.190:6379
>>> Performing Cluster Check (using node 10.254.0.190:6379)
M: ac5d899e8694a8a22fba55c817288b2b6d6b5029 10.254.0.190:6379
   slots:0-5460 (5461 slots) master
   1 additional replica(s)
M: 36508c8d0838c4978b4afb673737ccd0687177de 10.254.0.193:6379@16379
   slots: (0 slots) master
   0 additional replica(s)
S: 27e2faf486bd1e9a134b4303c1c5cc9de13fbb55 10.254.1.12:6379@16379
   slots: (0 slots) slave
   replicates 8a04fdb87d98c18ca37b6292965acf3adbbf0d1f
M: f0d5598e3dc10315576ee50b47b22a7bf0a03f5a 10.254.1.11:6379@16379
   slots:5461-10922 (5462 slots) master
   1 additional replica(s)
S: a0dd57c7dfb2331a399277868ae23a8cf279f839 10.254.0.191:6379@16379
   slots: (0 slots) slave
   replicates ac5d899e8694a8a22fba55c817288b2b6d6b5029
M: 8a04fdb87d98c18ca37b6292965acf3adbbf0d1f 10.254.2.11:6379@16379
   slots:10923-16383 (5461 slots) master
   1 additional replica(s)
S: 2d22d42fc9c8e44891eb4ebd8a5d1c757fbdf904 10.254.2.12:6379@16379
   slots: (0 slots) slave
   replicates f0d5598e3dc10315576ee50b47b22a7bf0a03f5a
[OK] All nodes agree about slots configuration.
>>> Check for open slots...
>>> Check slots coverage...
[OK] All 16384 slots covered.
>>> Send CLUSTER MEET to node 10.254.0.194:6379 to make it join the cluster.
[OK] New node added correctly.
```

还没完, 虽然现在`cluster info`中`cluster_known_nodes`值变成了8, 但是使用`cluster nodes`的结果中你会发现, 有两个节点是没有`slot`的(`slot`的值在`0-16383`之间).

```
127.0.0.1:6379> cluster nodes
6e826abcbbea6c97fa781e0cf3fdb822b5c1e0c0 10.254.0.200:6379@16379 master - 0 1584563291865 7 connected
859f18ac27c1772d0c8cd5449c04a68b457502f2 10.254.2.16:6379@16379 slave 43ea0ad2229250d27780ce19ebfac15720d796dc 0 1584563292783 6 connected
ad89d521c4d164efa78628963b2920bdac814235 10.254.1.16:6379@16379 slave ff117b911f60a59eaa7a47f11ad80db2027c0fe0 0 1584563292883 5 connected
43ea0ad2229250d27780ce19ebfac15720d796dc 10.254.1.15:6379@16379 myself,master - 0 1584563291000 3 connected 10923-16383
730e815556b1e0debbc0f825960fc8977f548f4d 10.254.0.199:6379@16379 master - 0 1584563293397 0 connected
b72bd81260b91f540a2c8e79de29ce8c0700290f 10.254.0.198:6379@16379 slave ed5201befc4d2c474b71ee4ddccec2575271f95b 0 1584563291865 4 connected
ff117b911f60a59eaa7a47f11ad80db2027c0fe0 10.254.2.15:6379@16379 master - 0 1584563291558 2 connected 5461-10922
ed5201befc4d2c474b71ee4ddccec2575271f95b 10.254.0.197:6379@16379 master - 0 1584563292579 1 connected 0-5460
```

> `10.254.0.199`和`10.254.0.200`是新增的两个节点, 虽然两个都成了master(现在成了`master * 5 + slave * 3`).

对比之前的slot状态.

```
127.0.0.1:6379> cluster nodes
859f18ac27c1772d0c8cd5449c04a68b457502f2 10.254.2.16:6379@16379 slave 43ea0ad2229250d27780ce19ebfac15720d796dc 0 1584559188772 6 connected
ad89d521c4d164efa78628963b2920bdac814235 10.254.1.16:6379@16379 slave ff117b911f60a59eaa7a47f11ad80db2027c0fe0 0 1584559189282 5 connected
43ea0ad2229250d27780ce19ebfac15720d796dc 10.254.1.15:6379@16379 myself,master - 0 1584559187000 3 connected 10923-16383
b72bd81260b91f540a2c8e79de29ce8c0700290f 10.254.0.198:6379@16379 slave ed5201befc4d2c474b71ee4ddccec2575271f95b 0 1584559189793 4 connected
ff117b911f60a59eaa7a47f11ad80db2027c0fe0 10.254.2.15:6379@16379 master - 0 1584559189000 2 connected 5461-10922
ed5201befc4d2c474b71ee4ddccec2575271f95b 10.254.0.197:6379@16379 master - 0 1584559187744 1 connected 0-5460
```

所以还需要重新划分`slot`.

```console
$ clusterNode=$(dig +short redis-app-0.redis-service.default.svc.cluster.local):6379
$ redis-trib reshard $clusterNode
How many slots do you want to move (from 1 to 16384)? 3000                 ## 要向指定节点迁移多少个slot(此时5个节点, 每个节点平均在3000多左右, 这里就填3000好了)
What is the receiving node ID? 730e815556b1e0debbc0f825960fc8977f548f4d    ## 这3000个slot要迁移到哪个node? 
Please enter all the source node IDs.
  Type 'all' to use all the nodes as source nodes for the hash slots.
  Type 'done' once you entered all the source nodes IDs.
Source node #1:all                                                         ## `all` 表示从所有节点上的slot抽调3000个slot到目标节点
    Moving slot 997 from ed5201befc4d2c474b71ee4ddccec2575271f95b
    Moving slot 998 from ed5201befc4d2c474b71ee4ddccec2575271f95b
   ...省略 这里打印了3000行, 就是要抽调的那3000个slot
Do you want to proceed with the proposed reshard plan (yes/no)? yes        ## 确认
Moving slot 997 from 10.254.0.197:6379 to 10.254.0.199:6379@16379:
Moving slot 998 from 10.254.0.197:6379 to 10.254.0.199:6379@16379:
...省略 这里是迁移过程
```

```
127.0.0.1:6379> cluster nodes
6e826abcbbea6c97fa781e0cf3fdb822b5c1e0c0 10.254.0.200:6379@16379 master - 0 1584566193244 7 connected
b72bd81260b91f540a2c8e79de29ce8c0700290f 10.254.0.198:6379@16379 slave ed5201befc4d2c474b71ee4ddccec2575271f95b 0 1584566192531 4 connected
859f18ac27c1772d0c8cd5449c04a68b457502f2 10.254.2.16:6379@16379 slave 43ea0ad2229250d27780ce19ebfac15720d796dc 0 1584566193450 6 connected
ff117b911f60a59eaa7a47f11ad80db2027c0fe0 10.254.2.15:6379@16379 myself,master - 0 1584566191000 2 connected 6462-10922
ad89d521c4d164efa78628963b2920bdac814235 10.254.1.16:6379@16379 slave ff117b911f60a59eaa7a47f11ad80db2027c0fe0 0 1584566193000 5 connected
730e815556b1e0debbc0f825960fc8977f548f4d 10.254.0.199:6379@16379 master - 0 1584566193971 8 connected 0-998 5461-6461 10923-11921
43ea0ad2229250d27780ce19ebfac15720d796dc 10.254.1.15:6379@16379 master - 0 1584566193451 3 connected 11922-16383
ed5201befc4d2c474b71ee4ddccec2575271f95b 10.254.0.197:6379@16379 master - 0 1584566193451 1 connected 999-5460
```

现在再看`730e815556b1e0debbc0f825960fc8977f548f4d`, 可以看到这上面的slot范围为: `0-998`, `5461-6461`, `10923-11921`.


