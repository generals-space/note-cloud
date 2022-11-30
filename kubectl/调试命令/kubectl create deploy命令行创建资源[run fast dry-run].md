# kubectl create deploy命令行创建资源

参考文章

1. [kubectl for Docker Users](https://kubernetes.io/docs/reference/kubectl/docker-cli-to-kubectl/)

快速创建一个 deploy 资源, 仅指定有限个字段, 无需
```
k create deploy mydeploy --image=registry.cn-hangzhou.aliyuncs.com/generals-space/centos:7:latest 
```

> k create 子命令只能创建 deploy, 无法创建 sts, ds 等同类型的资源(创建pod可以直接用`k run`子命令).

## 最大限制 command

作用有限, 这种方式没办法指定`command`, `--`也不起作用, 如下

```
k create deploy mydeploy --image=registry.cn-hangzhou.aliyuncs.com/generals-space/centos:7:latest -- tail -f /etc/os-release
```

kubectl set子命令也没有设置command的功能.

## --dry-run -oyaml 快速生成 deploy 模板

```
k create deploy mydeploy --image=registry.cn-hangzhou.aliyuncs.com/generals-space/centos:7:latest --dry-run -oyaml
```

`--dry-run`只创建对象但不提交到 apiserver.

```yaml
## k create deploy hjl-deploy --image=k8s-deploy/centos --dry-run -oyaml
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: hjl-deploy
  name: hjl-deploy
spec:
  replicas: 1
  selector:
    matchLabels:
      app: hjl-deploy
  strategy: {}
  template:
    metadata:
      creationTimestamp: null
      labels:
        app: hjl-deploy
    spec:
      containers:
      - image: k8s-deploy/centos
        name: centos
        resources: {}
status: {}
```
