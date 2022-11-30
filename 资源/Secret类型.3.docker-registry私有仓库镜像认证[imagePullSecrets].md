# Secret使用-私有仓库镜像认证[imagePullSecrets docker-registry]

参考文章

1. [kubernetes拉取私有仓库image方法](https://blog.csdn.net/kozazyh/article/details/79427119)
2. [Using kubernetes init containers on a private repo](https://stackoverflow.com/questions/42462244/using-kubernetes-init-containers-on-a-private-repo)
3. [kubernetes init containers using a private repo](https://stackoverflow.com/questions/42422892/kubernetes-init-containers-using-a-private-repo)
4. [官方文档 Menggunakan Init Container](https://kubernetes.io/id/docs/concepts/workloads/pods/init-containers/#menggunakan-init-container)
5. [How to use private repository with initContainer](https://stackoverflow.com/questions/53185465/how-to-use-private-repository-with-initcontainer)
6. [initContainers does not accept imagePullSecrets](https://github.com/kubernetes/kubernetes/issues/70732)

命令行可以通过`docker login 仓库地址(不必加https, 默认为docker.io)`登录, 之后可以保存类似session一样的会话, 之后就不必再登录了.

但是kubernetes是不行的, 但kubernetes针对这种有情况有专门的解决方法, 那就是`secret`资源对象.

创建示例

```
kubectl create secret docker-registry secret资源名称 --docker-server=仓库地址(只要域名, 无需路径) --docker-username=用户名 --docker-password=密码 --docker-email=xxx@gmail.com
```

- `--docker-server`: 仓库地址
- `--docker-username`: 仓库登陆账号
- `--docker-password`: 仓库登陆密码
- `--docker-email`: 邮件地址(必填)
- `-n`: 命名空间, 默认为当前命名空间

然后在kubernetes配置文件中, 指定使用的`secret`资源. 如下

```yaml
spec:
  imagePullSecrets:
  - name: secret资源名称
  containers:
  - name: test
    image: hub.c.163.com/xxx:latest
```

`imagePullSecrets`块的认证配置可以被所有`containers[]`与`initContainers[]`共用.

------

type为`kubernetes.io/dockerconfigjson`.

```console
$ kubectl create secret docker-registry registry-secret --docker-server=192.168.1.1 --docker-username=admin --docker-password=123456 --docker-email=xxx@gmail.com
secret/registry-secret created
$ k get secret
NAME               TYPE                              DATA   AGE
registry-secret    kubernetes.io/dockerconfigjson    1      5s
```

