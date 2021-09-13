# kafka pod重启后logstash无法获取新pod的IP

参考文章

1. [Client Not Following Broker IP Change - Upgrade Kafka Client to 2.1.1](https://github.com/akka/alpakka-kafka/issues/734)

问题描述

kube容器环境, logstash+kafka, 使用sts形式部署.

logstash中配置了kafka的 headless-service 地址.

```conf
input {
    kafka {
        bootstrap_servers => ["kafka-headless:9092"]
        ##...
    }
}
```

当 kafka pod 重启后, pod IP 发生了变动, 但是logstash一直使用原来的 pod IP, 然后无法拉取到日志信息, 除非重启 logstash pod, 才会重新获取正确的 pod IP.

------

听起来很像是 kube 内部的 dns 缓存问题, 当初我就是这么想的, 但是这个缓存存在的时间太长了, 不会过期.

后来发现是 Kafka 客户端(即 logstash)的 bug, 见参考文章1.

> This behavior exists in the **2.1.0** Kafka client. There is a release for Kafka **2.1.1** that resolves this issue.

