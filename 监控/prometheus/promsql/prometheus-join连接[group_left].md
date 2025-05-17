## 示例1

有如下两个指标, 想以`metallb_speaker_announced`为主, 以`node`字段将两者连接, 追加`node_network_info`的`address`和`device`字段, 两者的value则不关注, 适合放在Table面板中展示.

```sql
$ metallb_speaker_announced

metallb_speaker_announced{ip="192.168.82.19", job="metallb-speaker", node="master-1", service="default/test-service"}   1
metallb_speaker_announced{ip="192.168.82.45", job="metallb-speaker", node="master-1", service="metallb/webhook-svc"}    1
metallb_speaker_announced{ip="192.168.82.22", job="metallb-speaker", node="master-3", service="remote/remote-dev-svc"}  1
```

```sql
$ node_network_info{device="eth0"}

node_network_info{address="b4:83:51:19:12:62", device="eth0", job="node_exporter", node="master-1"}                     1
node_network_info{address="b4:83:51:19:0d:d6", device="eth0", job="node_exporter", node="master-2"}                     1
node_network_info{address="b4:83:51:19:12:64", device="eth0", job="node_exporter", node="master-3"}                     1
```

连接的效果如下

```sql
$ sum(metallb_speaker_announced{}) by (ip,node,service)
   * on(node) group_left(address,device)
(node_network_info{device="eth0"})

{address="b4:83:51:19:12:62", device="eth0", ip="192.168.82.19", node="white-master-1", service="default/test-service"} 1
{address="b4:83:51:19:12:62", device="eth0", ip="192.168.82.45", node="white-master-1", service="metallb/webhook-svc"}  1
{address="b4:83:51:19:12:64", device="eth0", ip="192.168.82.22", node="white-master-3", service="remote/remote-dev-svc"}1
```

`sum by(ip,node,service)`表示保留`metallb_speaker_announced`的这3个字段.

注意: `group_left`后面的指标还可以作为过滤条件, 实现用另一个指标去过滤主目标.
