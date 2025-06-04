# docker-容器日志配置

参考文章

1. [Docker官方文档 Configure logging drivers](https://docs.docker.com/config/containers/logging/configure/)

```json
{
  "log-driver": "json-file",
  "log-opts": {
    "max-size": "10m",
    "max-file": "3",
    "labels": "production_status",
    "env": "os,customer"
  }
}
```

上述代码为`/etc/docker/daemon.json`中日志相关的配置. 但日志驱动不只是在dockerd中配置, 还可以在启动容器时, 通过`docker run`的`--log-driver`, `--log-opts`命令指定, 格式相同.

docker支持多种日志驱动, 默认为`json-file`.

官方文档列举出了所有支持的日志驱动, 但是docker-ce版本只支持3种: local, json-file, journald.

除了`json-file`和`journal`, 其他的日志驱动无法使用`docker logs 容器id`查看到.

------

Docker 提供两种模式来将日志消息从容器发送到日志驱动程序:

1. (默认)直接阻塞式的从容器发送到驱动程序
2. 非阻塞发送, 将日志消息存储在中间每个容器的环形缓冲区中供驱动程序使用(...类似于生产者消费者呗)

非阻塞消息传送模式可防止应用程序因记录的压力(logging back pressure)而被阻塞. 当 STDERR 或 STDOUT 流阻塞时, 应用程序可能会以意想不到的方式失败.

> 警告：当缓冲区已满且新消息排入队列时, 内存中最早的消息将被丢弃. 丢弃消息通常首选阻止应用程序的日志写入过程. (Dropping messages is often preferred to blocking the log-writing process of an application.)

- `mode`: 这个日志选项用于控制使用阻塞还是非阻塞哪个消息发送方式. 
- `max-buffer-size`: 这个日志选项用于控制非阻塞方式下用作中间消息存储的环形缓冲区大小, 默认是 1MB. 

```
docker run -it --log-opt mode=non-blocking --log-opt max-buffer-size=4m alpine ping 127.0.0.1
```
