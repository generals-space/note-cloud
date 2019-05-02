# docker启动容器报错-Error response from daemon

参考文章

1. ["Error response from daemon: service endpoint with name es already exists" when starting container](https://github.com/moby/moby/issues/20398)

2. [docker: Error response from daemon: service endpoint with name XXX already exists.](http://blog.csdn.net/awewong/article/details/78516926)

```
$ docker run --name task-server --net=cpo --ip=172.18.1.12 -p 8202:8112 -d --rm -v /logs/cpo/task-server:/logs reg01.sky-mobi.com/cpo/task-server:1.0.0 --eureka.client.service-url.defaultZone=http://172.16.4.40:8761/eureka/,http://172.16.4.40:8762/eureka/ --spring.cloud.config.enabled=true --spring.cloud.config.discovery.enabled=true --spring.cloud.config.discovery.service-id=config-server --spring.cloud.config.username=user --spring.cloud.config.password=123456 --spring.profiles.active=test
41ef624105a50583343bcd9ba033f2d4f7beed9552c8c5f08c720ec1e3690717
docker: Error response from daemon: endpoint with name task-server already exists in network cpo.
```

场景描述: 

就是启动一个容器的时候, 就报错了. 最初以为是有停掉的容器依然占用着ip的缘故, 但是`docker ps -a`没有发现这个名为`task-server`的容器.

容器用的是自定义的网络, 名为`cpo`, 网段是`172.18`. 不过这和自定义网络应该没什么关系, 按照参考文章1中的说法, 应该是容器上次停止时返回码不正确, 导致网络中并没有真正移除它. 所以容器虽然不在了, 但是cpo网络中依然保留着这个容器所占用的IP, 不能被分配. 

```
$ docker network inspect cpo
[
    {
        "Name": "cpo",
        "Id": "36c0cf03a9c5ebb7a284972b9081c16e08b31f4c7e0404ceee7f681810faa7f0",
        "Scope": "local",
        "Driver": "bridge",

        "Containers": {
            "439a5c1904013c51aa71643a26771699a264d4b73fbe1b49c4bc2dd141ebb3ed": {
                "Name": "task-server",
                "EndpointID": "6dec7ba08b42034718ddc10ab52b5a1171f50ac08e602c8e9fd1f349485ccfb7",
                "MacAddress": "02:42:ac:12:01:0c",
                "IPv4Address": "172.18.1.12/16",
                "IPv6Address": ""
            },
            ...
        },
    }
]
```

按照参考文章2中的做法, 把`task-server`容器(实际已经不存在了)从cpo网络中强制卸载.

```
$ docker network disconnect -f cpo task-server
```

再次`inspect`网络时, `task-server`已经不在了, 启动容器成功.