# ctr导入导出docker镜像[namespace import export load save]

参考文章

1. [containerd 手动导入镜像](https://www.zeng.dev/post/2020-containerd-image-import/)
    - containerd 镜像有 namespace 的概念
2. [containerd 导入镜像](https://www.cnblogs.com/dream397/p/13815280.html)

kubernetes 使用 containerd 作 runtime 后, 引发一个问题: containerd 的镜像无法直接引用 docker 构建的镜像, 这两个是完全不互通的.

ctr 也拥有 pull 的命令, 所以对于远程仓库的镜像, 可以 containerd 直接拉取.

但当我使用 docker build 在本地构建了一个镜像时(并没有push), 删除当前主机上的 pod, 希望其引用的镜像同时更新, containerd 就找不到.

因此需要一个将 docker 镜像转换成 containerd 的镜像的方法, 主要是给 kubernetes 使用.

按照参考文章1, ctr 是有 namespaces 概念的, kubernetes 使用 containerd 下载的镜像, 都放在 k8s.io 这个 namespaces 下.

```
docker save k8s.gcr.io/pause -o pause.tar
ctr -n k8s.io images import pause.tar
```

导出

```
ctr -n k8s.io images export pause.tar pause:latest
```
