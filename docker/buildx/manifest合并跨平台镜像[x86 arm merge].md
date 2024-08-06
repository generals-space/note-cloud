# manifest合并跨平台镜像[x86 arm]

x86主机上也可以拉取arm镜像, 只是不能启动.

A  双平面地址
B  x86地址
C  arm地址

docker manifest rm A
docker manifest create A B C
docker manifest push A
