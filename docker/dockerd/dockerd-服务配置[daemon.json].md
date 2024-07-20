# dockerd-服务配置

参考文章

1. [官方文档 dockerd cli](https://docs.docker.com/engine/reference/commandline/dockerd/)
    - dockerd在命令行运行的可用选项及解释
2. [官方文档 Configure and troubleshoot the Docker daemon](https://docs.docker.com/config/daemon/)
    - dockerd的命令行选项与daemon.json配置文件选项的映射示例

```json
{
    "registry-mirrors": [
        "https://hub-mirror.c.163.com", 
        "https://registry.docker-cn.com", 
        "https://docker.mirrors.ustc.edu.cn"
    ],
    "dns": ["223.5.5.5", "223.6.6.6"],
    "max-concurrent-downloads" : 20,
    "max-concurrent-uploads" : 20
}
```

> docker-ce将`graph`字段修改为`data-root`.

其中`data-root`字段为docker所有的镜像, 容器存放的位置, 该目录不必预先存在, 启动docker服务时会自动创建.

