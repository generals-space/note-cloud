# kubernetes+postgres

## NFS挂载卷

参考文章

1. [chown: /var/lib/postgresql/data/postgresql.conf: Read-only file system](https://stackoverflow.com/questions/51884999/chown-var-lib-postgresql-data-postgresql-conf-read-only-file-system)

2. [How to deploy Postgresql on Kubernetes with NFS volume](https://stackoverflow.com/questions/51725559/how-to-deploy-postgresql-on-kubernetes-with-nfs-volume)

PV,PVC的配置这里不详细介绍, 踩坑本来是因为NFS. NFS对共享目录默认的权限设置为`no_all_squash,root_squash`

```yml
  template:
    metadata:
      labels:
        app: postgres
    spec:
      volumes:
        - name: pgnfs-vol
          persistentVolumeClaim:
            claimName: pgnfs-claim
      containers:
      - name: postgres
        image: postgres:9.6-alpine
        volumeMounts:
        - name: pgnfs-vol
          mountPath: "/var/lib/postgresql/data"
          readOnly: false 
```

但是在启动时报错

```
runuser: cannot set groups: operation not permitted
```

貌似是因为没有该目录的写权限? 但是为了配合`no_all_squash`, 我已经把共享目录的属主修改为了`nfsnobody`, 除非postgres容器需要使用root用户对这个目录进行一些修改, 按照参考文章1和2, 验证了我的猜想, 将`root_squash`修改为`no_root_squash`, 再次创建容器成功.

## 自定义启动命令

1. [“root” execution of the PostgreSQL server is not permitted](https://stackoverflow.com/questions/28311825/root-execution-of-the-postgresql-server-is-not-permitted)

有时需要修改数据库的最大连接数, 但我对postgres的配置文件并不十分熟悉, 按照postgres在dockerhub中对镜像的介绍, 可以通过`-c`选项指定字段去设置值.

```
docker run -d --name some-postgres postgres -c 'shared_buffers=256MB' -c 'max_connections=200'
```

但是在kubernetes配置时需要注意, 要在`command`中指定`postgres`, 否则无法找到执行文件.

```yml
      containers:
      - name: postgres
        image: postgres:9.6-alpine
        command: ["postgres", "-c", "max_connections=2000"]
```

但是这样在启动时出现如下报错.

```
"root" execution of the PostgreSQL server is not permitted.
The server must be started under an unprivileged user ID to prevent
possible system security compromise.  See the documentation for
more information on how to properly start the server.
```

我找了半天, 没有找到kubernetes配置中可以指定执行用户的方法. 后来在查看postgres的dockerfile时发现, ta的ENTRYPOINT设置为`docker-entrypoint.sh`, 所以可以通过这个脚本代为执行启动命令.

```yml
      containers:
      - name: postgres
        image: postgres:9.6-alpine
        - name: pgnfs-vol
          mountPath: "/var/lib/postgresql/data"
          readOnly: false 
        command: ["docker-entrypoint.sh", "-c", "max_connections=2000"]
```