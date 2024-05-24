参考文章

1. [openshift 4.5.9 etcd损坏+脑裂修复过程](https://zhangguanzhang.github.io/2021/06/08/ocp4.5.9-restore-etcd/#/%E6%93%8D%E4%BD%9C)
    - 跟我遇到的情况类似
    - snapshot save 卡住...
    - crictl没有cp命令, 生成的文件拷贝不出来...
    - 将生成的备份文件移动到共享目录, 发现etcd容器里面没有`mv`命令...
2. [How to backup etcd on a Kubernetes cluster created with kubeadm - rpc error: code = 13](https://stackoverflow.com/questions/51370870/how-to-backup-etcd-on-a-kubernetes-cluster-created-with-kubeadm-rpc-error-cod)

## 问题描述

进入到etcd容器内部执行如下命令进行备份, 但是一直卡住...

```log
$ ETCDCTL_API=3 etcdctl --endpoints=127.0.1:2379 snapshot save snapshotdb
{"level":"info","ts":"2023-03-10T02:14:08.257Z","caller":"snapshot/v3_snapshot.go:65","msg":"created temporary db file","path":"snapshotdb.part"}
```

添加上`--debug=true`查看日志.

```log
WARNING: 2023/03/10 02:11:37 [core] Adjusting keepalive ping interval to minimum period of 10s
WARNING: 2023/03/10 02:11:37 [core] Adjusting keepalive ping interval to minimum period of 10s
INFO: 2023/03/10 02:11:37 [core] parsed scheme: "etcd-endpoints"
INFO: 2023/03/10 02:11:37 [core] ccResolverWrapper: sending update to cc: {[{172.30.1.2:2379 172.30.1.2 <nil> 0 <nil>}] 0xc0004a8a60 <nil>}
INFO: 2023/03/10 02:11:37 [core] ClientConn switching balancer to "round_robin"
INFO: 2023/03/10 02:11:37 [core] Channel switches to new LB policy "round_robin"
INFO: 2023/03/10 02:11:37 [balancer] base.baseBalancer: got new ClientConn state:  {{[{172.30.1.2:2379 172.30.1.2 <nil> 0 <nil>}] 0xc0004a8a60 <nil>} <nil>}
{"level":"info","ts":"2023-03-10T02:11:37.370Z","caller":"snapshot/v3_snapshot.go:65","msg":"created temporary db file","path":"snapshotdb.part"}
INFO: 2023/03/10 02:11:37 [core] Subchannel Connectivity change to CONNECTING
INFO: 2023/03/10 02:11:37 [core] Subchannel picks a new address "172.30.1.2:2379" to connect
INFO: 2023/03/10 02:11:37 [balancer] base.baseBalancer: handle SubConn state change: 0xc000498a00, CONNECTING
INFO: 2023/03/10 02:11:37 [core] Channel Connectivity change to CONNECTING
INFO: 2023/03/10 02:11:37 [transport] transport: loopyWriter.run returning. connection error: desc = "transport is closing"
INFO: 2023/03/10 02:11:37 [core] Subchannel Connectivity change to TRANSIENT_FAILURE
INFO: 2023/03/10 02:11:37 [balancer] base.baseBalancer: handle SubConn state change: 0xc000498a00, TRANSIENT_FAILURE
INFO: 2023/03/10 02:11:37 [core] Channel Connectivity change to TRANSIENT_FAILURE
INFO: 2023/03/10 02:11:38 [core] Subchannel Connectivity change to IDLE
INFO: 2023/03/10 02:11:38 [balancer] base.baseBalancer: handle SubConn state change: 0xc000498a00, IDLE
INFO: 2023/03/10 02:11:38 [core] Subchannel Connectivity change to CONNECTING
INFO: 2023/03/10 02:11:38 [core] Subchannel picks a new address "172.30.1.2:2379" to connect
INFO: 2023/03/10 02:11:38 [balancer] base.baseBalancer: handle SubConn state change: 0xc000498a00, CONNECTING
INFO: 2023/03/10 02:11:38 [transport] transport: loopyWriter.run returning. connection error: desc = "transport is closing"
INFO: 2023/03/10 02:11:38 [core] Subchannel Connectivity change to TRANSIENT_FAILURE
```

## 解决方法

按照参考文章2中所说, 在容器里执行命令, 也是需要加上证书和key等参数的.

```log
$ etcdctl --endpoints=127.0.0.1:2379 --cacert=/etc/kubernetes/pki/etcd/ca.crt --cert=/etc/kubernetes/pki/etcd/server.crt --key /etc/kubernetes/pki/etcd/server.key snapshot save snapshotdb
{"level":"info","ts":"2023-03-10T02:14:08.257Z","caller":"snapshot/v3_snapshot.go:65","msg":"created temporary db file","path":"snapshotdb.part"}
{"level":"info","ts":"2023-03-10T02:14:08.276Z","logger":"client","caller":"v3@v3.5.6/maintenance.go:212","msg":"opened snapshot stream; downloading"}
{"level":"info","ts":"2023-03-10T02:14:08.276Z","caller":"snapshot/v3_snapshot.go:73","msg":"fetching snapshot","endpoint":"127.0.0.1:2379"}
{"level":"info","ts":"2023-03-10T02:14:08.669Z","logger":"client","caller":"v3@v3.5.6/maintenance.go:220","msg":"completed snapshot read; closing"}
{"level":"info","ts":"2023-03-10T02:14:08.678Z","caller":"snapshot/v3_snapshot.go:88","msg":"fetched snapshot","endpoint":"127.0.0.1:2379","size":"6.1 MB","took":"now"}
{"level":"info","ts":"2023-03-10T02:14:08.678Z","caller":"snapshot/v3_snapshot.go:97","msg":"saved","path":"snapshotdb"}
Snapshot saved at snapshotdb
```

成功.
