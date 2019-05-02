删除旧容器, 不会删除当前正在运行的容器

```
docker rm $(docker ps -a -q)
```

移除本地多余的历史镜像

```
docker images | grep '<none>' | awk '{print $3}' | xargs docker rmi
```
