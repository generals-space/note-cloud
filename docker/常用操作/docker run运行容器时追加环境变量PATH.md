# Docker run运行容器时追加环境变量PATH

参考文章

1. [How do I append to PATH environment variable when running a Docker container?](https://stackoverflow.com/questions/37281533/how-do-i-append-to-path-environment-variable-when-running-a-docker-container)

使用如下方式启动时, 容器内部的PATH值将会变成宿主机本身的PATH, 因为bash在运行时会将`$PATH`替换掉, 但这并不是我们想要的.

```
docker run -it -e PATH=$PATH:/usr/local/go/bin generals/centos7 /bin/bash
```

参考文章1提到了, 这个问题似乎无解, docker并没有更优雅的解决方法, 但也提供了一个折中方案.

```
docker run -it generals/centos7 bash -c 'export PATH=$PATH:/user/local/go/bin; bash'
```
