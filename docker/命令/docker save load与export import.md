# docker save load与export import

参考文章

1. [docker save与docker export的区别](https://cnodejs.org/topic/59a2304f7aeedce818249eeb)
    - `docker save`保存的是镜像(image), `docker export`保存的是容器(container); 
    - `docker load`用来载入镜像包, `docker import`用来载入容器包, **但两者都会恢复为镜像**; 
    - `docker load`不能对载入的镜像重命名, 而`docker import`可以为镜像指定新名称;
    -
    - `docker export`的应用场景主要用来制作基础镜像, 比如你从一个ubuntu镜像启动一个容器, 然后安装一些软件和进行一些设置后, 使用docker export保存为一个基础镜像, 然后把这个镜像分发给其他人使用.

`save`与`load`的常用命令.

```bash
## busybox 为镜像名称
docker save -o busybox.tar busybox
docker load < ./busybox.tar
```

```bash
## nginx 为容器名称(可以是 stop 的容器)
docker export -o nginx.tar nginx
```

`docker save`保存的tar包, 解压开是`manifest.json`, `repositories`和一堆hash目录, 其中存储着该镜像的所有分层.

`docker export`保存的tar包, 解压开是`bin`, `boot`, `dev`, `etc`, `home`等和`.dockerenv`, 可以看出是一个OS的完整目录.
