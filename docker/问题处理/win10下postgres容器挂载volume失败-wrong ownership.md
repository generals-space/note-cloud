# win10下postgres容器挂载volume失败-wrong ownership

参考文章

1. [Mounting data volume for Postgres in docker for Windows doesn't work](https://github.com/docker/for-win/issues/445)

2. [Data directory “/var/lib/postgresql/data/pgdata” has wrong ownership](https://forums.docker.com/t/data-directory-var-lib-postgresql-data-pgdata-has-wrong-ownership/17963/31)

在win10下创建postgres容器, 希望能够挂载目录防止数据丢失, 按照如下命令启动(`testing-pgdata`目录需要事先存在), 但是会出错退出.

```
docker run -d --name testing-pg -v d:\dockerdata\testing-pgdata:/var/lib/postgresql/data postgres:11
```

使用`docker logs`查看容器日志发现有如下报错.

```
2019-04-27 03:59:53.364 UTC [79] FATAL:  data directory "/var/lib/postgresql/data" has wrong ownership
2019-04-27 03:59:53.364 UTC [79] HINT:  The server must be started by the user that owns the data directory.
child process exited with exit code 1
```

网上查到的消息最终都指向`docker-forwin`项目本身的bug(见参考文章1), 而且貌似至今未解决. 

后来按照参考文章2中所说, 在windows下挂载目录, 需要使用`docker volume create`先创建一个挂载卷, 然后在启动命令中使用这个卷, 就可以了, 如下.

```
## 创建挂载卷
$ docker volume create testing-pgdata -d local
testing-pgdata
## 查看已经存在的挂载卷
$ docker volume ls
DRIVER              VOLUME NAME
local               testing-pgdata
## 启动容器
$ docker run -d --name testing-pg -v testing-pgdata:/var/lib/postgresql/data postgres:11
```

这样启动就可以了.

...只不过这样创建的挂载卷没有明确的路径, 我们没有办法知道目录在宿主机上的实际位置. 至少目前我还没找到.

------

但是, 如果你使用`docker volume create`创建了一个挂载卷, 而你想在`docker-compose.yml`中使用时?

```yml
version: '3'
services:
  postgres-svc:
    image: postgres:11
    ports:
    - "5432:5432"
    ## win下挂载路径可能会出错
    volumes:
      ## - ./data/postgres:/var/lib/postgresql/data
      - testing-pgdata:/var/lib/postgresql/data
```

```
$ docker-compose up -d 
ERROR: Named volume "testing-pgdata:/var/lib/postgresql/data:rw" is used in service "postgres-svc" but no declaration was found in the volumes section.
```

这样的操作是不对的, 正确的方法如下, 不需要事先手动创建挂载目录, 而是直接将其声明到yml文件中.

```yml
version: '3'
services:
  postgres-svc:
    image: postgres:11
    ports:
    - "5432:5432"
    ## win下挂载路径可能会出错
    volumes:
      ## - ./data/postgres:/var/lib/postgresql/data
      - testing-pgdata:/var/lib/postgresql/data
volumes:
  testing-pgdata:
```

这样在启动compose时就会自动创建挂载卷.

```
$ docker-compose up -d postgres-svc
Creating volume "bootstrapmb-downloader-pyasync_testing-pgdata" with default driver
Creating bootstrapmb-downloader-pyasync_postgres-svc_1 ... done
```
