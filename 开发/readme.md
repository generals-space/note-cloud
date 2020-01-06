## 20191230 `client-go/tools/leaderelection/resourcelock`: 分布式资源锁. 

参考文章

1. [谈谈k8s的leader选举--分布式资源锁](https://blog.csdn.net/weixin_39961559/article/details/81877056)
    - 作用讲解得很清楚, 给出了应用场景
2. [kubernetes 自定义控制器的高可用](http://blog.fatedier.com/2019/04/17/k8s-custom-controller-high-available/)
    - 深入代码, 原理分析

用于多副本资源的同步行为, 竞争成功的才能进行操作. `scheduler`和`controller-manager`都用到了这个机制, `etcd-operator`中也是, `backup controller`也用到了此机制, 锁竞争成功的副本才可以创建CRD, 然后进行其他的行为.

