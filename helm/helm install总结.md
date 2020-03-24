## 安装指定版本的chart

安装

```
helm install kubeapps [-n 命名空间] bitnami/kubeapps
```

安装指定版本

```
helm install kubeapps [-n 命名空间] bitnami/kubeapps --version 1.2.0
```

更新

```
helm upgrade kubeapps [-n 命名空间] -f 目标values文件 bitnami/kubeapps
```

```
$ helm uninstall kubeapps
release "kubeapps" uninstalled
```

使用`helm fetch`下来chart工程并解压, 查看`values.yaml`哪些可以修改的字段, 写到`myval.yaml`文件中, 然后使用如下命令安装

```
helm install kubeapps -f myval.yaml bitnami/kubeapps
```
