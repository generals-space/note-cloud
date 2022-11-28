# kuber-拉取私有仓库镜像

参考文章

1. [kubernetes拉取私有仓库image方法](https://blog.csdn.net/kozazyh/article/details/79427119)
2. [Using kubernetes init containers on a private repo](https://stackoverflow.com/questions/42462244/using-kubernetes-init-containers-on-a-private-repo)
3. [kubernetes init containers using a private repo](https://stackoverflow.com/questions/42422892/kubernetes-init-containers-using-a-private-repo)
4. [官方文档 Menggunakan Init Container](https://kubernetes.io/id/docs/concepts/workloads/pods/init-containers/#menggunakan-init-container)
5. [How to use private repository with initContainer](https://stackoverflow.com/questions/53185465/how-to-use-private-repository-with-initcontainer)
6. [initContainers does not accept imagePullSecrets](https://github.com/kubernetes/kubernetes/issues/70732)

命令行可以通过`docker login 仓库地址(不必加https, 默认为docker.io)`登录, 之后可以保存类似session一样的会话, 之后就不必再登录了.

但是kubernetes是不行的, 但kubernetes针对这种有情况有专门的解决方法, 那就是`secret`资源对象.

`secret`类型目前有3种: 

1. docker-registry(仓库配置)
2. generic(通用类型, 任意文件/目录/键值对)
3. tls(密钥配置)

创建`secret`资源

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

## initContainers 部分的 image 如何拉取私有镜像?

kuber版本: 1.16.2

flannel的yaml部署文件包括 `initContainers` 和 `containers` 两部分, `imagePullSecrets`定义的密钥能共用吗?

我创建了阿里云镜像仓库的镜像, 但是拉取失败了, 且就是在 `initContainers` 部分失败的(用 `describe` 命令查看得知). 

本来以为是 kuber 不支持, 找到了参考文章2, 2引用了3, 3引用了4, 说其中的解决方法没用. 后来我又找到了5, 引用了官方 issue 即参考文章6, 当时最后一个 issue 说人家1.15版本的明明实验成功了啊, 这个 issue 为啥还不关...

然后回头看参考文章2, 提问者最后补充说明, `imagePullSecrets` 其实能被两者共用, 我回头一看, 创建 secret 的命令, 用户名和仓库地址写反了...

但这种情况是两者属于同一私有仓库的, 如果这两者引用的是不同私有仓库的镜像呢(虽然这样好像不太合理啊)? ~~总感觉还是分开搞比较好~~. 我觉得还是不要考虑这个问题了, 因为就算是`containers`部分, 容器也不是只有一个, 要是还考虑这些不同容器来自不同的私有仓库未免将问题过于复杂化了.

安安生生干活吧.
