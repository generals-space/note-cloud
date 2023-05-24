# 官方文档参考手册

参考文章

1. [Assign Static IP to Docker Container and Docker-Compose](https://www.baeldung.com/ops/docker-assign-static-ip-container)
    - docker 指定静态IP
        1. `docker network create --subnet=10.5.0.0/16 mynet`
        2. `docker run --net mynet --ip 10.5.0.1 -e MYSQL_ROOT_PASSWORD=123456 mysql:latest`
    - docker compose 指定静态IP: driver, ipam, subnet
2. [jmarcos-cano/docker-compose.redis.yml](https://gist.github.com/jmarcos-cano/9b5d16d5e65999875e60f44be4c370df)
    - docker 指定 sysctl 内核参数(本文其实并没有纯 docker 的实现, 这里只是留作对比和备忘)
        1. `docker run --sysctl kernel.sem="250 6400000 1000 25600" -e MYSQL_ROOT_PASSWORD=123456 mysql:latest`
        2. `docker run --sysctl kernel.sem=net.core.somaxconn=1024 -e MYSQL_ROOT_PASSWORD=123456 mysql:latest`
    - docker compose 指定内核参数
3. [docker_sysctl](https://beixiu.net/dev/docker-sysctl/)
    - 同参考文章2
4. [Set secomp to unconfined in docker-compose](https://stackoverflow.com/questions/46053672/set-secomp-to-unconfined-in-docker-compose)
    - docker 指定`security-opt`
        1. `docker run --security-opt seccomp=unconfined -e MYSQL_ROOT_PASSWORD=123456 mysql:latest`
    - docker compose 指定`security-opt`
5. [openGauss / openGauss-container](https://gitee.com/opengauss/openGauss-container/tree/03dd83be9dbd19b4ee302bd4f379a6ce5c1d4f1b/)
    - opengauss 官方仓库

以下配置创建的服务为 opengauss 集群(1主2从), 不过无所谓, 只要看 yaml 配置即可, 对应的 docker 启动命令可见参考文章5.

```yaml
version: '3'
services:
  opengauss-01:
    hostname: primary
    image: opengauss:5.0.0
    networks:
      opengauss-net:
        ipv4_address: 172.12.0.2
    volumes:
    - /data/opengauss_volume:/volume
    environment:
    - primaryhost=172.12.0.2
    - primaryname=primary
    - standbyhosts=172.12.0.3,172.12.0.4
    - standbynames=standby1,standby2
    - GS_PASSWORD=test@123
    ulimits:
      nofile:
        soft: "640000"
        hard: "640000"
    sysctls:
      kernel.sem: "250 6400000 1000 25600"
    security_opt: 
    - seccomp:unconfined
  opengauss-02:
    hostname: primary
    image: opengauss:5.0.0
    networks:
      opengauss-net:
        ipv4_address: 172.12.0.3
    volumes:
    - /data/opengauss_volume:/volume
    environment:
    - primaryhost=172.12.0.2
    - primaryname=primary
    - standbyhosts=172.12.0.3,172.12.0.4
    - standbynames=standby1,standby2
    - GS_PASSWORD=test@123
    ulimits:
      nofile:
        soft: "640000"
        hard: "640000"
    sysctls:
      kernel.sem: "250 6400000 1000 25600"
    security_opt: 
    - seccomp:unconfined
  opengauss-03:
    hostname: primary
    image: opengauss:5.0.0
    networks:
      opengauss-net:
        ipv4_address: 172.12.0.4
    volumes:
    - /data/opengauss_volume:/volume
    environment:
    - primaryhost=172.12.0.2
    - primaryname=primary
    - standbyhosts=172.12.0.3,172.12.0.4
    - standbynames=standby1,standby2
    - GS_PASSWORD=test@123
    ulimits:
      nofile:
        soft: "640000"
        hard: "640000"
    sysctls:
      kernel.sem: "250 6400000 1000 25600"
    security_opt: 
    - seccomp:unconfined

networks:
  opengauss-net:
    driver: bridge
    ipam:
      config:
      - subnet: 172.12.0.0/24
        ## 参考文章1中的配置有 gateway , 但是在实践时该字段会报错.
        ## gateway: 172.12.0.1
```
