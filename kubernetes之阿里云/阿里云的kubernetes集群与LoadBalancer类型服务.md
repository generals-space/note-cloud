# 阿里云的kubernetes集群LoadBalancer类型service对象与SLB对象

<!--
<!date!>: 2018-12-25
-->

参考文章

1. [阿里云Serverless Kubernetes通过Ingress提供7层服务访问](https://www.jianshu.com/p/7da30664b84a)

2. [通过负载均衡（Server Load Balancer）访问服务](https://help.aliyun.com/document_detail/86531.html)

标准kubernetes概念中service的类型有

1. `ClusterIP`

2. `LoadBalancer`

在一个`service`配置中, 如果设置了`type`为`LoadBalancer`, 阿里云会在创建此`service`的同时, 还会创建一个公网的SLB, 且公网IP随机.

之后在容器控制台中**服务模块**对应的service资源就会出现刚刚建好的SLB的公网地址.

![](https://gitee.com/generals-space/gitimg/raw/master/7b441b410557e6f3d2f51d43850e8c23.png)

我们希望每个service资源都使用同一个已经存在的SLB实例, 而不是每建一个`LoadBalancer`都新建一个公网SLB, 不好管理还费钱!

参考文章1中**使用说明**部分说明了阿里云kubernetes集群中SLB的生成机制. 并且参考文章1和参考文章2都讲了可以通过`annotations`指定`LoadBalancer`类型的service对象所绑定的SLB. 如下

```yml
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