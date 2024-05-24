# 由metrics-server引发的apiserver无法启动的问题[APIService]

参考文章

1. [ks-apiserver 无法启动](https://kubesphere.com.cn/forum/d/3017-ks-apiserver)
    - k8s的服务发现机制要求所有apiservice都是True状态
2. [kubernetes问题解答专栏](https://www.ziji.work/kubernetes/kubernetes-question-answer-special.html)
3. [kubernetes 1.18.x metrics-server 采集不到数据](https://www.haxi.cc/archives/k8s-1-18-x-metrics-server-no-data.html)

## 问题描述

kubernetes: v1.17.3

master-0: 172.22.248.225
master-1: 172.22.248.226
master-2: 172.22.248.227
worker-0: 172.22.248.227

某天早上上班后, 发现kubectl获取集群信息失败了, 显示如下报错

```
The connection to the server 172.22.248.239:6443 was refused - did you specify the right
```

`172.22.248.239`是master前面的负载均衡服务器, 上去看了看, ta本身没什么问题. 于是将248.225上的`admin.conf`中, 集群的地址从`248.239`改成`248.225`, 然后重新发起请求.

```
export KUBECONFIG=/etc/kubernetes/admin.conf
```

但是仍然报错.

用`docker ps | grep apiserver`可以看到, apiserver是启动了的, 只不过由于主机宕机, 发生过重启, 但是基本上所有组件都恢复正常了.

docker logs 查看日志, 发现有如下输出

![](https://gitee.com/generals-space/gitimg/raw/master/e11b626d2d063b3e4d2935a9d1f7ca7c.jpg)

```
E0826 04:25:37.762890       1 controller.go:114] loading OpenAPI spec for "v1beta1.metrics.k8s.io" failed with: failed to retrieve openAPI spec, http error: ResponseCode: 503, Body: service unavailable, Header: map[Content-Type:[text/plain; charset=utf-8] X-Content-Type-Options:[nosniff]]
```

在这个集群里有部署`metrics-server`, 不过已经挂了很长时间了, 没人用也没人修. 但是`metrics-server`是一个第三方组件, 理论上应该不影响核心组件的正常运行, 所以一开始我没有把目光放到`metrics-server`服务上. 但是根据`failed to retrieve openAPI spec, http error`作为关键字搜索时, 几乎全都与`metrics-server`有关系...

经过一番查找后, 找到了参考文章1, 其中有一句特别重要: **k8s的服务发现机制要求所有apiservice都是True状态**.

这让我怀疑起了`CrashBackOff`的`metrics-server`会不会真的影响到了`apiserver`的启动流程.

我将`admin.conf`中连接的集群地址从225改成226, 又从226改成了227, 终于发现连接227时, kubectl是可以正常请求的, 就是说227的apiserver还是正常的.

`metrics-server`的Pod异常, 所以Endpoint也是异常, 所以ta对应的APIService资源是False状态.

```log
$ kubectl get apiservice | grep metrics
NAME                        SERVICE                        AVAILABLE                 AGE
v1beta1.metrics.k8s.io      kube-system/metrics-server     False (MissingEndpoints)  15m
```

ok, 现在通过227删除`metrics-server`的Pod是解决不了问题的, 只能先将`v1beta1.metrics.k8s.io`这个`apiservice`资源删掉了.

完成之后, 225, 226上的apiserver都可以访问了, 239也可以正常转发了, 问题解决(注意备份).

------

225, 226上的apiserver都挂了, 但是为什么227上的可以访问?

我的猜测是, 如果apiserver本身没有选主机制的话, 那就是之前3个apiserver是正常启动的状态, 之后`metrics-server`挂掉, 225, 226主机先后宕机, 只剩下227一个在服务.

在ta们三个主机上执行`docker ps | grep apiserver`查看启动时间, 227上的apiserver启动时间是最久的, 印证了我的猜测.

这么说来, 坏掉的APIService可真是危险啊...
