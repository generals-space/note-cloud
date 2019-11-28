# 允许非root用户使用docker命令

参考文章

1. [Manage Docker as a non-root user](https://docs.docker.com/install/linux/linux-postinstall/#manage-docker-as-a-non-root-user)

docker在安装时就创建了一个命名docker的用户组, 只要将指定用户添加到docker组中, 就能使用docker命令来管理集群.

```
cat /etc/group | grep docker
usermod -aG docker 指定用户
```
