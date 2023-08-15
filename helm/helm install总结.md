## 安装指定版本的chart

### 安装

```
helm install kubeapps [-n 命名空间] bitnami/kubeapps
```

### 安装指定版本

```
helm install kubeapps [-n 命名空间] bitnami/kubeapps --version 1.2.0
```

### 更新

```
helm upgrade kubeapps [-n 命名空间] -f 目标values文件 bitnami/kubeapps
```

```
$ helm uninstall kubeapps
release "kubeapps" uninstalled
```
