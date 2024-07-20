部署kubernetes指定containerd时, 需要做一些准备工作.

首先生成默认配置文件

```
containerd config default > /etc/containerd/config.toml
```

对配置文件做如下修改.

```yaml
## 1. 注释 disabled_plugins 字段
## disabled_plugins = ["cri"]

## 2. 修改 pause 镜像地址
[plugins]
  [plugins."io.containerd.grpc.v1.cri"]
    sandbox_image = "registry.cn-hangzhou.aliyuncs.com/google_containers/pause:3.6"
```

重启containerd服务

```
systemctl restart containerd
```

用如下命令长久设置 endpoint

```
crictl config runtime-endpoint unix:///run/containerd/containerd.sock
```
