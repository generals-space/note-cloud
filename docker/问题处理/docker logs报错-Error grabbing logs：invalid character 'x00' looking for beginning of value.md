# docker logs报错-Error grabbing logs：invalid character 'x00' looking for beginning of value

参考文章

1. [Error grabbing logs: invalid character '\x00' looking for beginning of value #140](https://github.com/docker/for-linux/issues/140)

## 问题描述

`docker logs -f`查看实时日志, 但是直接结束退出了.

```log
root@k8s-master-01:~# docker logs -f django
[I 15:41:43.357           runserver:168] HTTP GET /api/system/auth/namespace 200 [0.03, 172.21.0.3:51438]
[I 15:42:10.793           runserver:168] HTTP GET /api/system/auth/namespace 200 [0.04, 172.21.0.3:52936]
Error grabbing logs: invalid character '\x00' looking for beginning of value
```

容器本身没有问题, 多次docker logs可以看到日志在持续刷新, 只是没办法自动刷.

## 原因分析

参考文章1中有人回复, 可以用grep找到日志文件中的非法字符, 删掉就行.

```log
grep -P '\x00' /var/lib/docker/containers/**/*json.log
```

> log路径可以通过`docker inspect django | grep log`查看.

但是日志是在实时写入的, 我本来担心手动修改日志文件会引发错乱, 后来又看到一个回复.

```
docker logs -f django --tail=20
```

直接用`--tail`参数, 跳过有乱码的行就可以了.
