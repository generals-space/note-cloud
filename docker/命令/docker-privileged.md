`docker run`时不使用`--privileged`选项启动容器, 可以在`docker exec`时指定, 也可以使用超级权限的命令.

使用`--privileged`启动的容器, 重启后依然可以保留`pviledge`权限.
