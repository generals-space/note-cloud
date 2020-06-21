参考文章

1. [一文快速上手Logstash](https://cloud.tencent.com/developer/article/1353068)

2. [logstash收集TCP端口日志](https://www.cnblogs.com/Dev0ps/p/9314551.html)

## 认识

在ELK系统里, E与K已经十分固定, 没什么可定制的空间. 业内对日志收集层面的调整和优化倒是层出不穷, filebeat+kafka前置于logstash, 也可以直接将日志收集到elesticsearch. 所以关于ELK的使用实例主要集中在日志收集的方式展示上, 当然还有kibana的日志查询和聚合语句的使用.

## logstash配置解释

`./logstash/pipeline`目录下有几个处理文件.

### `stdio.conf`

读取tcp 5000端口的消息作为日志, 应该用到了tcp插件. 这也是docker-compose中为logstash服务映射5000端口的原因: 方便在宿主机上进行测试. 向这个商品发送的消息会在logstash容器的标准输出中打印. 

```
$ telnet logstash 5000
Trying 172.24.0.5...
Connected to logstash.
Escape character is '^]'.
hello world
good for u
^]
telnet> quit
Connection closed.
```

```
{
    "@timestamp" => 2019-07-24T03:31:32.339Z,
       "message" => "hello world\r",
      "@version" => "1",
          "host" => "elk_nginx_1.elk_default",
          "port" => 35246
}
{
    "@timestamp" => 2019-07-24T03:31:35.493Z,
       "message" => "good for u\r",
      "@version" => "1",
          "host" => "elk_nginx_1.elk_default",
          "port" => 35246
}
```
