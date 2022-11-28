# 阿里云的kuber集群与LoadBalancer Service

<!--
<!date!>: 2018-12-25
-->

参考文章

1. [阿里云Serverless Kubernetes通过Ingress提供7层服务访问](https://www.jianshu.com/p/7da30664b84a)

2. [通过负载均衡（Server Load Balancer）访问服务](https://help.aliyun.com/document_detail/86531.html)

标准kubernetes概念中service的类型有

1. `ClusterIP`: 只能在集群内部访问(通过service名称), 如果要暴露给公网访问, 需要配合使用Ingress对象.
2. `NodePort`: 将会在kube-proxy服务所在的节点上(基本是所有节点...Master和Worker)"创建"一个30000-32767端口, 访问节点上的此端口就相当于访问pod内部的服务. 
    - nodePort不一定真正创建, 工作在在iptables模式的kube-proxy服务使用转发链完成此操作, 用netstat/ss是看不到listening状态的端口的.
3. `LoadBalancer`: 这种类型的Service需要云厂商支持, 自建集群无法创建, 因为没有办法为此service实例赋予externalIP.
    - 实际上, 像GCE, 阿里云, ta们的LoadBalancer就是绑定了SLB实例的NodePort类型的service而已, 由公有云的SLB代理service开放给公网.

在一个`service`配置中, 如果设置了`type`为`LoadBalancer`, 阿里云会在创建此`service`的同时, 还会创建一个公网的SLB, 且公网IP随机.

之后在容器控制台中**服务模块**对应的service资源就会出现刚刚建好的SLB的公网地址.

![](https://gitee.com/generals-space/gitimg/raw/master/7b441b410557e6f3d2f51d43850e8c23.png)

我们希望每个service资源都使用同一个已经存在的SLB实例, 而不是每建一个`LoadBalancer`都新建一个公网SLB, 不好管理还费钱!

参考文章1中**使用说明**部分说明了阿里云kubernetes集群中SLB的生成机制. 并且参考文章1和参考文章2都讲了可以通过`annotations`指定`LoadBalancer`类型的service对象所绑定的SLB. 如下

```yaml
apiVersion: v1
kind: Service
metadata:
  annotations:
    service.beta.kubernetes.io/alicloud-loadbalancer-id: "your_loadbalancer_id"
  name: nginx
  namespace: default
spec:
  ports:
  - port: 443
    protocol: TCP
    targetPort: 443
  selector:
    run: nginx
  type: LoadBalancer
```

不知道纯kubernetes集群中是否有`annotations`的概念来指定SLB, 因为没有吧.

注意: 

1. 预先使用`LoadBalancer`创建service资源时, 如果不指定SLB id, 自动生成的SLB是没有办法解绑的. 所以只能在创建时指定好一个已经存在的SLB. 如果已经自动生成了, 而且生成了两个, 如上面图中就生成了两个, 就只能把两个service的`type`都改成`Cluster`, 删除自动生成的SLB, 才能解绑, 之后才能绑定我们想要绑定的SLB.
