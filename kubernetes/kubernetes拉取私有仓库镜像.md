# kubernetes拉取私有仓库镜像

参考文章

1. [kubernetes拉取私有仓库image方法](https://blog.csdn.net/kozazyh/article/details/79427119)

命令行可以通过`docker login 仓库地址(不必加https, 默认为docker.io)`登录, 之后可以保存类似session一样的会话, 之后就不必再登录了.

但是kubernetes是不行的, 但kubernetes针对这种有情况有专门的解决方法, 那就是`secret`资源对象.

`secret`类型目前有3种: 

1. docker-registry(仓库配置)

2. generic(通用类型, 任意文件/目录/键值对)

3. tls(密钥配置)

创建`secret`资源

```
kubectl create secret docker-registry secret资源名称 --docker-server=仓库地址 --docker-username=用户名 --docker-password=密码 --docker-email=xxx@gmail.com
```

- `--docker-server`: 仓库地址

- `--docker-username`: 仓库登陆账号

- `--docker-password`: 仓库登陆密码

- `--docker-email`: 邮件地址(必填)

- `-n` 命名空间, 默认为当前命名空间

然后在kubernetes配置文件中, 指定使用的`secret`资源. 如下

```yml
spec:
    imagePullSecrets:
    - name: secret资源名称
    containers:
    - name: test
    image: hub.c.163.com/xxx:latest
```