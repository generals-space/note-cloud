为了体现 sts 在分布式集群, 有状态服务中的作用, 这里我们使用 deployment 资源来部署 es.

kibana 也是 deployment.

logstash 用于收集各主机上的日志, 采集的目标就是 nginx, 所以这两个需要使用 daemonset 类型进行部署.

按照这个结构部署的 ELK, 就算已经在 yaml 中写明了用户名密码分别为`elastic/123456` 访问 kibana webUI 是不需要密码的, 访问 es 的 http 接口也不需要使用密码.

```console
$ curl es:9200/_cat/health
1592740248 11:50:48 elasticsearch green 1 1 2 2 0 0 0 0 - 100.0%
```

## 配置文件后缀必须为 .yml

如果命名为 .yaml, 容器将无法启动(`CrashLoopBackOff`状态), 查看日志有如下报错.

```console
$ k logs -f fb17c160244c
Exception in thread "main" SettingsException[elasticsearch.yaml was deprecated in 5.5.0 and must be renamed to elasticsearch.yml]
        at org.elasticsearch.node.InternalSettingsPreparer.prepareEnvironment(InternalSettingsPreparer.java:72)
        at org.elasticsearch.cli.EnvironmentAwareCommand.createEnv(EnvironmentAwareCommand.java:95)
        at org.elasticsearch.cli.EnvironmentAwareCommand.execute(EnvironmentAwareCommand.java:86)
        at org.elasticsearch.cli.Command.mainWithoutErrorHandling(Command.java:124)
        at org.elasticsearch.cli.MultiCommand.execute(MultiCommand.java:77)
        at org.elasticsearch.cli.Command.mainWithoutErrorHandling(Command.java:124)
        at org.elasticsearch.cli.Command.main(Command.java:90)
        at org.elasticsearch.common.settings.KeyStoreCli.main(KeyStoreCli.java:41)
```

## 密码 123456

es 本身可以将密码设置为123456(yaml里用双引号包裹即可), 但是 kibana 不行, 就算有双引号, 也会启动失败, 有如下日志.

```
 FATAL  Error: [elasticsearch.password]: expected value of type [string] but got [number]
```

...不过貌似在配置文件中写 123456 就没问题...难道说只是不能在环境变量里写?
