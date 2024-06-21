# coredns PodReady endpoints与域名解析

参考文章

1. [kafka broker 注册到 ip 上](https://blog.csdn.net/weixin_42195284/article/details/91042571)

## 场景描述

部署一个 kafka 集群, 拥有 statefulset 与 headless service 资源.

statefulset 中是通过对 kafka 监听的 tcp 端口的检测作为健康检查依据的.

```yaml
    readinessProbe:
      failureThreshold: 3
      tcpSocket:
        port: 9092
      periodSeconds: 10
      successThreshold: 1
      timeoutSeconds: 10
```

kafka 需要注册到 zk 做服务发现, 但是各实例所在 Pod 用来注册的地址写成了自己的域名地址, 日志如下

```log
Registered broker 0 at path /brokers/ids/0 with addresses: PLAINTEXT -> EndPoint(test-01-0.test-01-svc.kube-system.svc.cluster.local,9092,PLAINTEXT) (kafka.utils.ZkUtils)
```

Pod-0、Pod-1、Pod-2 需要进行通信, 就只能通过这个地址来.

但这就导致了一个问题◑ˍ◐

1. 各 kafka 实例需要建立集群, 选主后才算正式启动, 然后9092端口才可访问;
2. 9092端口的健康检查通不过, 不仅影响将 podIP 注册到 endpoints, 还影响 coredns 对其访问域名的解析;
    - coredns不需要解析 not ready 状态的 pod, 访问上述域名直接就会显示: `Name or service not know`
3. 域名访问不了就导致各 kafka 实例没办法相互通信, 然后不算启动完成, 9092端口不通.

于是就出现了死循环, 导致集群一直就起不来.

但是在该集群之外, 连接ta所属的 zk, 你会发现`/brokers/ids`路径下已经拥有了3个id.

## 解决方案

1. 通过添加配置项`advertised.listeners=PLAINTEXT://192.168.16.11:9092`, 让kafka注册到zk的地址改为IP(默认是hostname域名)
2. 编写健康检查脚本, 在容器中连接zk, 将zk中注册的kafka实例数量与期望数量对比, 作为判断依据

