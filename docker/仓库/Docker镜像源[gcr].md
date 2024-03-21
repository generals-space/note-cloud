参考文章

1. [k8s官方镜像代理加速](https://blog.csdn.net/chengyinwu/article/details/134843205)
2. [无法拉取 gcr.io 镜像？用魔法来打败魔法](https://mp.weixin.qq.com/s/Vt0FRTx1PsoYFdLa0QZzWw)
3. [无法愉快拉取 gcr.io、quay.io、ghcr.io 容器镜像？手把手教你用魔法来打败魔法](https://blog.csdn.net/easylife206/article/details/124763001)
    - 参考文章2的转载文章

gcr.io -> gcr.dockerproxy.com

registry.k8s.io -> k8s.mirror.nju.edu.cn

## 常见的镜像仓库

docker.io: DockerHub 的官方仓库, 也是 Docker 的默认仓库

gcr.io, k8s.gcr.io: 谷歌镜像仓库

quay.io: RedHat 镜像仓库

ghcr.io: GitHub 镜像仓库

## 常见的国内镜像源

中国区官方镜像: https://registry.docker-cn.com
清华源: https://docker.mirrors.ustc.edu.cn
阿里源: https://cr.console.aliyun.com
腾讯源: https://mirror.ccs.tencentyun.com
网易源: http://hub-mirror.c.163.com
