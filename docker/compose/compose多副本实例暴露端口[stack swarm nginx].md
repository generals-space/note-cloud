# compose多副本实例暴露端口[stack swarm]

参考文章

1. [linux运维 4-compose/stack/swarm集群](https://www.kancloud.cn/noahs/linux/1116267)
2. [docker compose scale的时候解决端口冲突问题](https://blog.csdn.net/wangmarkqi/article/details/101155592)
3. [How to Use Docker Compose to Run Multiple Instances of a Service in Development](https://pspdfkit.com/blog/2018/how-to-use-docker-compose-to-run-multiple-instances-of-a-service-in-development/)

## deploy.replicas 多副本

docker-compose V3之后,设置资源限额和容器数量的配置,只能写在deploy字段里, 但是`docker-compose up` 却不支持 deploy 配置

编辑同一份docker-compose.yml, 但是compose 和 swarm/stack的分工是这样的:

- docker-compose: 用于`dev`支持build/restart, 但是不支持`deploy`;
- swarm/stack: 用于`prod`支持`deploy`的各种设置, 包括分配cpu和内存, 但是创建容器只支持从`image`, 不支持`build`.

本来想用docker-compose创建一个redis cluster集群的, 在yaml文件如下

```yaml
version: '3'
services:
  mcp-redis:
    container_name: mcp-redis
    image: registry.cn-hangzhou.aliyuncs.com/generals-space/redis:5.0.8.1
    restart: always
    deploy:
      mode: replicated
      replicas: 6
    networks:
    - mcp
    ports:
    - 6379:6379
networks:
  mcp:
    driver: bridge
```

但在执行up命令时报错

```console
$ dc up -d mcp-redis
WARNING: Some services (mcp-redis) use the 'deploy' key, which will be ignored. Compose does not support 'deploy' configuration - use `docker stack deploy` to deploy to a swarm.
Creating network "middleware_mcp" with driver "bridge"
Creating mcp-redis ... done
```

直接把`deploy`字段忽略了, 只创建了一个容器.

其实docker-compose已经过时了, 以后可以使用`docker compose up -d`.

## ports端口暴露

按照参考文章2, 3, 此问题基本无解, 想要为多个实例暴露同一个端口, 只能在前端配置一个nginx做负载均衡.
