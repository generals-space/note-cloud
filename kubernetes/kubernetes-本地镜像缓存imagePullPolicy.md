# kubernetes-本地镜像缓存imagePullPolicy

参考文章

1. [Kubernetes通过yaml配置文件创建实例时不使用本地镜像的原因](https://www.58jb.com/html/154.html)

```yml
spec: 
  containers: 
    - name: nginx 
      image: image: reg.docker.lc/share/nginx:latest 
      imagePullPolicy: IfNotPresent
```

`imagePullPolicy`可选值: 

1. `IfNotPresent`: 默认, 如果本地已经有该镜像, 不再重新pull.

2. `Never`: 直接不再去拉取镜像了, 使用本地的; 如果本地不存在就报异常了.

3. `Always`: 每次都尝试从远程pull镜像.