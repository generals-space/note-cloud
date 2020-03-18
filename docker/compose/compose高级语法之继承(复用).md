# compose高级语法之继承(复用)

参考文章

1. [consul 官方compose集群配置示例](https://github.com/hashicorp/consul/blob/master/demo/docker-compose-cluster/docker-compose.yml)

```yml
version: '3'

services:

  consul-agent-1: &consul-agent
    image: consul:latest
    networks:
      - consul-demo
    command: "agent -retry-join consul-server-bootstrap -client 0.0.0.0"

  consul-agent-2:
    <<: *consul-agent

  consul-agent-3:
    <<: *consul-agent

  consul-server-1: &consul-server
    <<: *consul-agent
    command: "agent -server -retry-join consul-server-bootstrap -client 0.0.0.0"

  consul-server-2:
    <<: *consul-server

  consul-server-bootstrap:
    <<: *consul-agent
    ports:
      - "8400:8400"
      - "8500:8500"
      - "8600:8600"
      - "8600:8600/udp"
    command: "agent -server -bootstrap-expect 3 -ui -client 0.0.0.0"

networks:
  consul-demo:
```

`&consul-agent`语句将`consul-agent-1`服务的配置声明为`consul-agent`, 之后的服务配置中可以使用`<<: *consul-agent`直接复用所有配置.

而如果出现与父级配置有出入的地方, 可以声明新的字段(如上面在`<<`又重新声明了`command`)来完成覆写.