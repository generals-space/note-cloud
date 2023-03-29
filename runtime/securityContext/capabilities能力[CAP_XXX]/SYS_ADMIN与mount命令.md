# SYS_ADMIN与mount命令

参考文章

1. [How do I mount --bind inside a Docker container?](https://stackoverflow.com/questions/36553617/how-do-i-mount-bind-inside-a-docker-container)

在容器里执行`mount`命令挂载一个目录, 需要容器拥有`CAP_SYS_ADMIN`能力.
