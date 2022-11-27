# kuber-nodePort, targetPort, port的区别和意义

参考文章

1. [Kubernetes中的nodePort，targetPort，port的区别和意义](https://blog.csdn.net/u013760355/article/details/70162242)

2. [从外部访问Kubernetes中的Pod](https://jimmysong.io/posts/accessing-kubernetes-pods-from-outside-of-the-cluster/)

3. [hostPort不生效](http://liupeng0518.github.io/2018/12/29/k8s/Network/%E5%BC%82%E5%B8%B8%E6%8E%92%E9%94%99/)
    - hostPort不生效的原因分析, 及使用iptables手动映射的解决方法

- `targetPort`: 表示该service要映射的源端口, 比如一个容器里监听的是80端口, 那`targetPort`就是80. 
- `port`: 是service监听的端口, pod中的服务相互访问时, 就是访问的其他pod绑定的service的端口.
- `nodePort`: kuber集群负责的, 开放给集群外部访问的端口, 范围在30000-32767. 由`kube-proxy`服务处理, 由于实际上集群中每个节点(包括master和worker)都运行着这个服务, 所以在每个节点上访问这个端口, 都能访问到ta对应的pod中的服务. 

kube-proxy服务工作在iptables模式下时, 节点上并没有真正监听nodePort(用netstat/ss是查不到的), 应该是使用iptables的转发链完成的. 至于ipvs模式下是否有监听, 还没有实验过.

------

还有一个`hostPort`, 在Pod, Deployment和DaemonSet对象的container字段中都可以使用.

```yaml
spec:
  containers:
  - name: nginx
    image: nginx
    ports:
    - name: nginx
      containerPort: 80
      hostPort: 80
```

ta的作用类似于使用`docker run`时的`-p`选项, 将容器的`containerPort`, 映射到pod被调度到的node节点的`hostPort`端口上.

由于pod可能被调度到不同节点上, 所以其实很不可靠, 除非限制pod被调度到的节点范围, 这一点用label和selector不难做到.

可能是由于网络插件的限制, 网上有很多说hostPort不生效的, 我也碰到了. 参考文章3提出了使用iptables手动解决映射的方法, 我想应该是有效的. 有一种更简单的方法, 就是使用docker的host网络一样, 让pod直接使用其所在节点的网络栈. 这需要用到`hostNetwork`字段.

```yaml
spec:
  hostNetwork: true
  containers:
  - name: nginx
    image: nginx
    ports:
    - name: nginx
      containerPort: 80
      hostPort: 80
```

由于容器直接使用宿主机网络, 此时`containerPort`与`hostPort`必须是一致的, 否则会出现如下错误.

```
The Pod "mypod" is invalid: spec.containers[0].ports[0].containerPort: Invalid value: 80: must match `hostPort` when `hostNetwork` is true
```