# prometheus容器启动失败-opening storage failed：data-chunks_head：invalid magic number 0

参考文章

1. [tsdb.Open fails with invalid magic number 0 when running with reverted previously mmaped chunks#7397](https://github.com/prometheus/prometheus/issues/7397#issuecomment-1175661042)
    - 官方issue

prometheus: v2.32.1

## 问题描述

在一次宿主机重启后, prometheus容器就无法启动了.

```
2ae27c85695c  prometheus:v2.32.1  "/opt/bitnami/promet…"  3 hours ago  Restarting (1) 10 seconds ago  prometheus
```

docker logs 查看容器日志.

```log
ts=2025-01-08T08:50:47.924Z caller=tls_config.go:195 level=info component=web msg="TLS is disabled." http2=false
ts=2025-01-08T08:50:47.928Z caller=main.go:799 level=info msg="Stopping scrape discovery manager..."
ts=2025-01-08T08:50:47.928Z caller=main.go:813 level=info msg="Stopping notify discovery manager..."
ts=2025-01-08T08:50:47.928Z caller=main.go:835 level=info msg="Stopping scrape manager..."
ts=2025-01-08T08:50:47.928Z caller=main.go:809 level=info msg="Notify discovery manager stopped"
ts=2025-01-08T08:50:47.928Z caller=manager.go:945 level=info component="rule manager" msg="Stopping rule manager..."
ts=2025-01-08T08:50:47.928Z caller=manager.go:955 level=info component="rule manager" msg="Rule manager stopped"
ts=2025-01-08T08:50:47.928Z caller=notifier.go:600 level=info component=notifier msg="Stopping notification manager..."
ts=2025-01-08T08:50:47.928Z caller=main.go:795 level=info msg="Scrape discovery manager stopped"
ts=2025-01-08T08:50:47.928Z caller=main.go:1055 level=info msg="Notifier manager stopped"
ts=2025-01-08T08:50:47.928Z caller=main.go:829 level=info msg="Scrape manager stopped"
ts=2025-01-08T08:50:47.928Z caller=main.go:1064 level=error err="opening storage failed: data/chunks_head/000344: invalid magic number 0"
```

容器的数据目录是挂载到宿主机的, 因此可以查看到出问题的文件

```log
root@localhost:/prometheus/chunks_head# ls -al
total 230296
drwxr-xr-x  2 root root      4096 Jan  8 16:30 .
drwxr-xr-x 35 root root      4096 Jan  8 16:51 ..
-rw-r--r--  1 root root 101459593 Jan  8 15:00 000342
-rw-r--r--  1 root root 134217587 Jan  8 16:30 000343
-rw-r--r--  1 root root    131072 Jan  8 16:30 000344
```

## 解决方法

参考文章1中说, 可以直接删除出问题的 chunks 文件, 不过可能会丢失部分数据.

我试了一下, 容器重启成功, 不过看监控数据好像也没丢什么数据, 可能是因为挂的时间比较短?(挂了24m)