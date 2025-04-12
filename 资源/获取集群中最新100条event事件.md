# 获取集群中最新100条event事件

## 场景描述

k8s默认只保存1小时内的event事件, 1小时前的会被自动清除.

当集群出现异常, event可能出现激增, 在超过1w条event时, kubectl get event会很耗费时间. 目前存在一个python系统, 通过 exec 方式执行 kubectl 查询 event 并返回, 展示到 web 界面, 查询超时默认为5s, 在这种情况下无法给出有效的响应, 因此考虑是否可以通过 kubectl 只查询最新100条事件.

但是k8s官方并没有提供查询最新100条event事件的方法, 只能通过 head/tail 命令做过滤, 但这是在客户端进行截取, 对操作耗时没有任何帮助.

```
kubectl get events -A --sort-by=.lastTimestamp | tail -n 500
```

因此考虑用 watch 长连接, 在本地新增一个缓存服务.

## kubectl -w + docker logs

kubectl get 没有`--tail`选项, 但是 kubectl logs 与 docker logs 是有`--tail`选项的, 考虑新增一个容器维持长连接, 然后在 python 服务中执行 docker logs --tail=100 查询最新事件.

```yaml
version: "3"
services:
  event:
    image: django:v3
    container_name: event
    entrypoint: 
    - bash
    - -c
    - >
      for (( ; ; ));
      do
        if ! $(kubectl get ns > /dev/null 2>&1); then
          sleep 15
          continue
        fi
        for i in {1..100}; do echo; done
        kubectl get event -A -o=jsonpath='{}{"\n"}' -w
      done
    environment:
    - KUBECONFIG=/usr/local/projects/admin.conf
    volumes:
    - ./myproject/admin.conf:/usr/local/projects/admin.conf # 集群的证书文件
    - /usr/bin/kubectl:/usr/bin/kubectl
    restart: always
    tty: true
    stdin_open: true
    deploy:
      resources:
        limits:
          memory: 100m
```

由于docker会保留容器日志, 而每次容器异常重启, kubectl get event -w 重新执行, 就会再获取1次全量事件. 为了避免python服务通过 docker logs --tail 获取到重复的事件, 需要在正常启动后, 打印100行空白行, 然后在 python 中对空行做判断.

这样存在一个问题, kubectl get event -w, 在获取全量事件列表后, 每新增一个事件都会打印出来.

但是同一个事件可能会重复上报, 每上报一次, count字段就会加一, 但是apiserver端会对其进行合并, kubectl get event 只展示一条, 而 kubectl get event -w 却会重复打印.

```log
LAST SEEN TYPE    REASON    OBJECT        SUBOBJECT                SOURCE                    MESSAGE                   FIRST SEEN   COUNT   NAME
108s      Warning Unhealthy pod/zero-db-0 spec.containers{zero-db} kubelet, green-master-2   Readiness probe failed:   156m         9422    zero-db-0.182ff49add27c2a5
```

此路不通.

## 新增 event-collector 服务

用golang编写 list-watch 服务, 并提供 http 接口, 每次通过`Lister.List()`然后根据lastTimestamp字段排序再截取, 速度很快.
