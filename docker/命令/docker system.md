# docker system

docker 1.13 时添加此命令.

system 有如下子命令

- `df`:     查看磁盘使用率
- `prune`:  方便的清理命令, 可以清理未运行的容器, 未被引用的网络, 无标签的镜像, 构建镜像时用到的缓存等...本机环境很有用.
- `events`: 与`docker events`完全一致
- `info`:   与`docker info`完全一致

## df

```log
$ docker system df
TYPE                TOTAL               ACTIVE              SIZE                RECLAIMABLE
Images              22                  4                   6.738GB             6.186GB (91%)
Containers          4                   0                   343B                343B (100%)
Local Volumes       40                  1                   652.8MB             652.8MB (99%)
Build Cache         0                   0                   0B                  0B
```

## prune

```log
$ docker system prune
WARNING! This will remove:
        - all stopped containers
        - all networks not used by at least one container
        - all dangling images
        - all dangling build cache
Are you sure you want to continue? [y/N] y
Deleted Containers:
4fcb5112a11366d4365c2388108687f4d9a16d5f55a61d3d5d771666fecc0d00
...省略

Deleted Networks:
compose_default
gowp_default
pycms_default
wuhougit_default

Deleted Images:
deleted: sha256:4dcae2a16d01a58e32d5762eb3a6a83f2c08999a303a6ce17f4255db7d01b984
...省略
Total reclaimed space: 24.85MB
```
