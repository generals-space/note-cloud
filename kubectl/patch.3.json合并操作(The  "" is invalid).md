# patch.3.json合并操作(The  "" is invalid)

参考文章

1. [kubectl patch增加或修改环境变量](https://blog.csdn.net/m0_37549390/article/details/118371216)

假设存在如下statefulset资源.

```yaml
spec:
  selector:
    matchLabels:
      middleware: logstash
      cluster: logstash-0321-01
  serviceName: logstash-0321-01-svc
  template:
    metadata:
      labels:
        middleware: logstash
        cluster: logstash-0321-01
    spec:
      affinity:
        nodeAffinity:
          requiredDuringSchedulingIgnoredDuringExecution:
            nodeSelectorTerms:
            - matchExpressions:
              - key: logstashSchedulable
                operator: NotIn
                values:
                - "false"

```

我们想在nodeAffinity中再添加一条规则, 如下

```yaml
              - key: mcp.logstash.log/pool
                operator: Exists
```

那么这个语句要怎么写?

```
kubectl patch sts logstash-0321-01 --type json -p '[{"op":"add","path":"/spec/template/spec/affinity/nodeAffinity/requiredDuringSchedulingIgnoredDuringExecution/nodeSelectorTerms/0/matchExpressions/-","value":{"key":"mcp.logstash.log/pool","operator":"Exists"}}]'
```

上述语句中`matchExpressions/-`后面那个`-`, 表示在`matchExpressions`数组后面追加一个成员. 这个值可以是一个确定的索引, 如`1`

```
kubectl patch sts logstash-0321-01 --type json -p '[{"op":"add","path":"/spec/template/spec/affinity/nodeAffinity/requiredDuringSchedulingIgnoredDuringExecution/nodeSelectorTerms/0/matchExpressions/1","value":{"key":"mcp.logstash.log/pool","operator":"Exists"}}]'
```

表示将新成员追加到第1个成员(索引是从0开始的), 但是这个值不能超过1, 否则会报如下错误

```console
$ kubectl patch sts logstash-0321-01 --type json -p '[{"op":"add","path":"/spec/template/spec/affinity/nodeAffinity/requiredDuringSchedulingIgnoredDuringExecution/nodeSelectorTerms/0/matchExpressions/2","value":{"key":"mcp.logstash.log/pool","operator":"Exists"}}]'
The  "" is invalid
```

## The  "" is invalid

由于向一个超过索引的位置插入新成员都会报这个错, 所以patch语句也需要根据情况进行调整, 如果statefulset的结构如下所示

```yaml
spec:
  selector:
    matchLabels:
      middleware: logstash
      cluster: logstash-0321-01
  serviceName: logstash-0321-01-svc
  template:
    metadata:
      labels:
        middleware: logstash
        cluster: logstash-0321-01
    spec:
      affinity:
        ## 避免同一个集群的不同pod调度到同一台主机
        podAntiAffinity:
          requiredDuringSchedulingIgnoredDuringExecution:
          - labelSelector:
              matchExpressions:
              - key: cluster
                operator: In
                values:
                - logstash-0321-01
            topologyKey: kubernetes.io/hostname
```

patch语句应该怎么写?

直接使用原来的语句, 会报错, 因为没有`nodeAffinity`规则, 所以其实是相当于越界了.

```
$ kubectl patch sts logstash-0321-01 --type json -p '[{"op":"add","path":"/spec/template/spec/affinity/nodeAffinity/requiredDuringSchedulingIgnoredDuringExecution/nodeSelectorTerms/0/matchExpressions/-","value":{"key":"mcp.logstash.log/pool","operator":"Exists"}}]'
The  "" is invalid
```

下面的语句会直接把"podAntiAffinity"换成"nodeAffinity"🤨, 并不是我们想要的.

```console
$ kubectl patch sts logstash-0321-01 --type json -p '[{"op":"add","path":"/spec/template/spec/affinity", "value":{"nodeAffinity": {"requiredDuringSchedulingIgnoredDuringExecution": {"nodeSelectorTerms": [{"matchExpressions": [{"key": "mcp.logstash.log/pool","operator": "Exists"}]}]}}}}]'
```

要用下面的

```console
$ kubectl patch sts logstash-0321-01 --type json -p '[{"op":"add","path":"/spec/template/spec/affinity/nodeAffinity", "value": {"requiredDuringSchedulingIgnoredDuringExecution": {"nodeSelectorTerms": [{"matchExpressions": [{"key": "mcp.logstash.log/pool","operator": "Exists"}]}]}}}]'
```
