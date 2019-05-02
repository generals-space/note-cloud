# Kubernetes-service类型ClusterIP, LoadBalancer和Ingress详解

参考文章

1. [Kubernetes的三种外部访问方式：NodePort、LoadBalancer和Ingress](https://mp.weixin.qq.com/s/2Rmca-kCoRp0TtHhuDtNNg)

2. [Kubernetes Ingress实战](http://www.cnblogs.com/zhaojiankai/p/7896357.html)

`targetPort`: 表示service要映射的源端口, 比如一个容器里监听的是80端口, 那`targetPort`就是80. 
