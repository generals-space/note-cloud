# calicoctl-Calico process is not running

参考文章

1. [calicoctl as a pod: Calico process is not running](https://github.com/projectcalico/calicoctl/issues/1594)

```log
$ calicoctl node status
Calico process is not running.
```

按照参考文章1中的说法, `calicoctl`在执行的时候会检测`bird`, `bird6`进程. 

看了下`calicoctl`的官方部署文件, 只写了`hostNetwork: true`, 尝试将添加`hostPID: true`, 然后就可以了.

但又出现了如下错误

```log
$ calicoctl node status
Calico process is running.

IPv4 BGP status
Error querying BIRD: unable to connect to BIRDv4 socket: dial unix /var/run/bird/bird.ctl: connect: no such file or directory
IPv6 BGP status
Error querying BIRD: unable to connect to BIRDv6 socket: dial unix /var/run/bird/bird6.ctl: connect: no such file or directory
```

使用ps查看bird进程

```log
$ ps -ef | grep bird
 8750 root      0:00 runsv bird
 8751 root      0:00 runsv bird6
 8873 root      2:06 bird -R -s /var/run/calico/bird.ctl -d -c /etc/calico/confd/config/bird.cfg
 8874 root      1:56 bird6 -R -s /var/run/calico/bird6.ctl -d -c /etc/calico/confd/config/bird6.cfg
```

看来还需要映射这个目录.

```yaml
  volumes:
    - name: vol-bird
      hostPath:
        path: /var/run/calico
        type: DirectoryOrCreate
  containers:
    - name: calicoctl
      ## ...省略
      volumeMounts:
        - name: vol-bird
          ## 挂载bird目录
          mountPath: /var/run/bird
```

重新部署后可正常执行

```log
$ calicoctl node status
Calico process is running.

IPv4 BGP status
+--------------+-------------------+-------+------------+-------------+
| PEER ADDRESS |     PEER TYPE     | STATE |   SINCE    |    INFO     |
+--------------+-------------------+-------+------------+-------------+
| 172.32.0.101 | node-to-node mesh | up    | 2019-12-10 | Established |
| 172.32.0.105 | node-to-node mesh | up    | 2019-12-10 | Established |
+--------------+-------------------+-------+------------+-------------+

IPv6 BGP status
No IPv6 peers found.
```

> `calico node status`显示是的各邻居节点的信息, 由于此时`calicoctl`容器被调度在`172.32.0.104`, 所以上述结果没有显示当前所在节点.

