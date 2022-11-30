# kubectl run创建临时pod

这个手段类似于 docker run, 创建一个临时pod, 用来验证网络通不通, 能否通过yum安装上某些工具包等.

```console
kubectl run mypod -it --image=registry.cn-hangzhou.aliyuncs.com/generals-space/centos:7:latest /bin/bash
kubectl run mypod -it --image=registry.cn-hangzhou.aliyuncs.com/generals-space/centos:7:latest -- /bin/bash
```

> `--`好像不是必要的, 等报错的时候再加也可以.

还可以加一些其他参数:

- `--rm`: 退出终端就删除
- `--restart=Never`: 自动重启
- `--port=8080 --expose`: 创建并绑定 service
    - `--expose`必须要有`--port`参数, 否则创建会失败
    - 单纯的`--port`参数只会在 pod yaml 中加上`containerPort`字段, 基本没什么作用.
- `--requests='cpu=100m,memory=256Mi'`
- `--limits='cpu=100m,memory=256Mi'`

## --dry-run -oyaml 快速生成 pod 模板

`--dry-run`只创建对象但不提交到 apiserver.

```yaml
## k run hjl-pod --image=k8s-deploy/centos --requests='cpu=100m,memory=256Mi' --dry-run -oyaml
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    run: hjl-pod
  name: hjl-pod
spec:
  replicas: 1
  selector:
    matchLabels:
      run: hjl-pod
  strategy: {}
  template:
    metadata:
      labels:
        run: hjl-pod
    spec:
      containers:
      - image: k8s-deploy/centos
        name: hjl-pod
        resources:
          requests:
            cpu: 100m
            memory: 256Mi
status: {}
```

------

kube 1.13版本中, `kubectl run`是通过`deployment`创建pod的;

kube 1.22版本中, `kubectl run`则是直接创建pod资源;
