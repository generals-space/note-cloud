```json
{
    "data-root": "/opt/docker",
    "registry-mirrors": [
        "https://registry.docker-cn.com", 
        "https://docker.mirrors.ustc.edu.cn"
    ]
}
```

> docker-ce将`graph`字段修改为`data-root`.

其中`data-root`字段为docker所有的镜像, 容器存放的位置, 该目录不必预先存在, 启动docker服务时会自动创建.

