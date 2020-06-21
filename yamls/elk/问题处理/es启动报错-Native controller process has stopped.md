# es启动报错-Native controller process has stopped
参考文章

1. [Elasticsearch修改network后启动失败](https://www.cnblogs.com/phpper/p/9803934.html)
2. [elasticsearch启动报错： Native controller process has stopped - no new native processes can be started](https://blog.csdn.net/K_Lily/article/details/105320221)

```
{"type": "server", "timestamp": "2020-06-21T09:17:24,362+0000", "level": "INFO", "component": "o.e.x.m.p.NativeController", "cluster.name": "elasticsearch", "node.name": "es-01",  "message": "Native controller process has stopped - no new native processes can be started"  }
```

实际上, 只在 initContainers 中加一句`sysctl -w vm.max_map_count=655300`就可以了.
