# kubectl工具配置

参考文章

1. [配置远程工具访问kubernetes集群](http://blog.csdn.net/shenshouer/article/details/52960364)

2. [Authenticate Across Clusters with kubeconfig](https://kubernetes.io/docs/tasks/access-application-cluster/authenticate-across-clusters-kubeconfig/)

当前我们通过`kubectl`工具与`apiserver`服务进行通信时, 都是通过`-s`指定目标`apiserver`的访问地址, 要写一长串. 尤其是连接不同的kuber集群时, 不同的`apiserver`地址会让人疯掉.

我们可以通过配置文件, 以别名的形式指定要连接的地址. 

`kubectl`可以通过`--kubeconfig`指定使用的配置文件的路径, 如果没有指定时, 默认使用`~/.kube/config`文件. 使用`kubectl config 子命令`可以修改这个文件中的内容.

方法如下

如果`~/.kube/config`文件不存在, 下面这条指令会自动创建.

```
$ kubectl config set-cluster sky-test --server=http://172.32.100.71:8080
Cluster "sky-test" set.
```

但是还是不能直接使用, 是为了方便多集群管理. 

```
$ kubectl get nodes
The connection to the server localhost:8080 was refused - did you specify the right host or port?
```

我们可以通过设置`context`表明我们要管理哪一个集群.

```
$ kubectl config set-context sky-test --cluster=sky-test
Context "sky-test" created.
$ kubectl config use-context sky-test
Switched to context "sky-test".
$ kubectl get nodes
NAME                    STATUS     AGE       VERSION
172.32.100.81           Ready       6d        v1.7.1-beta.0.2+09955ec93bcfc1
172.32.100.91           Ready       6d        v1.7.1-beta.0.2+09955ec93bcfc1
localhost.localdomain   Ready       6d        v1.7.1-beta.0.2+09955ec93bcfc1
```

看一下`~/.kube/config`的内容.

```yml
apiVersion: v1
clusters:          ## 定义集群  
- cluster:
    server: http://172.32.100.71:8080
  name: sky-test
contexts:           ## 定义上下文与哪一个集群相关联
- context:
    cluster: sky-test
    user: ""
  name: sky-test
current-context: sky-test   ## 定义当前正在管理的上下文.
kind: Config
preferences: {}
users: []
```