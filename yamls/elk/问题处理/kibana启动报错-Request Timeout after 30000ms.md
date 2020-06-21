# kibana启动报错-Request Timeout after 30000ms

参考文章

1. [kibana报错Request Timeout after 30000ms故障解决](https://blog.csdn.net/qq_40907977/article/details/104499178)
    - es资源设置太小导致 kibana 请求超时
2. [elasticsearch 7.2.0安装](https://blog.csdn.net/u011311291/article/details/100041912)
    - Elasticsearch cluster did not respond with license information.
3. [解决es集群启动完成后报master_not_discovered_exception](https://blog.csdn.net/qq_20042935/article/details/105274464)
    - 只启动了一个 es 节点, 按道理它就是 master 节点

集群 setup 完成后, 发现 kibana 总是陷入`CrashLoopBackOff`, 一直重启, 查看日志有如下输出.

```console
$ k logs -f kibana-XXXXXXXXX
{"type":"log","@timestamp":"2020-06-21T11:08:02Z","tags":["fatal","root"],"pid":1,"message":"{ Error: Request Timeout after 30000ms\n    at /usr/share/kibana/node_modules/elasticsearch/src/lib/transport.js:362:15\n    at Timeout.<anonymous> (/usr/share/kibana/node_modul
es/elasticsearch/src/lib/transport.js:391:7)\n    at ontimeout (timers.js:436:11)\n    at tryOnTimeout (timers.js:300:5)\n    at listOnTimeout (timers.js:263:5)\n    at Timer.processTimers (timers.js:223:10)\n  status: undefined,\n  displayName: 'RequestTimeout',\n  mes
sage: 'Request Timeout after 30000ms',\n  body: undefined,\n  isBoom: true,\n  isServer: true,\n  data: null,\n  output:\n   { statusCode: 503,\n     payload:\n      { statusCode: 503,\n        error: 'Service Unavailable',\n        message: 'Request Timeout after 30000
ms' },\n     headers: {} },\n  reformat: [Function],\n  [Symbol(SavedObjectsClientErrorCode)]: 'SavedObjectsClient/esUnavailable' }"}

 FATAL  Error: Request Timeout after 30000ms

```

但这个日志比较笼统, 不能准确定位问题. 比如参考文章1就说可能是由于 es 资源设置太小, 导致查询处理较慢, 但尝试修改了也没用.

通过 web 界面访问 kibana, 发现一片空白, 有输出如下信息.

```
Kibana server is not ready yet
```

后来我又回去翻了翻 kibana 的日志, 发现有很多像下面这样重复的信息.

```
{"type":"log","@timestamp":"2020-06-21T11:39:25Z","tags":["status","plugin:rollup@7.2.0","error"],"pid":1,"state":"red","message":"Status changed from yellow to red - [data] Elasticsearch cluster did not respond with license information.","prevState":"yellow","prevMsg":
"Waiting for Elasticsearch"}
{"type":"log","@timestamp":"2020-06-21T11:39:25Z","tags":["status","plugin:remote_clusters@7.2.0","error"],"pid":1,"state":"red","message":"Status changed from yellow to red - [data] Elasticsearch cluster did not respond with license information.","prevState":"yellow","
prevMsg":"Waiting for Elasticsearch"}
```

参考文章2中说需要打开`node.name`的注释(所以跟什么 license 完全没关系), 我联想到`node.name`与`cluster.initial_master_nodes`的关系. 我之前认为后者(数组)中的成员应该是节点的通信地址, 所以写的是 es 的 service 信息.

```yaml
node.name: es-01
## 这里的数组成员为 service 名称
cluster.initial_master_nodes: ["es"]
```

于是找了一个测试容器, 使用 curl 发送请求, 真的出错了.

```console
## es 是 service 名称
$ curl es:9200/_cat/health
{"error":{"root_cause":[{"type":"master_not_discovered_exception","reason":null}],"type":"master_not_discovered_exception","reason":null},"status":503}
```

这个报错就精确了很多.

于是把`cluster.initial_master_nodes`的成员改成了`node.name`的值, 重启 es, 接口可以正常访问, kibana也正常了.

```console
$ curl es:9200/_cat/health
1592740248 11:50:48 elasticsearch green 1 1 2 2 0 0 0 0 - 100.0%
```
