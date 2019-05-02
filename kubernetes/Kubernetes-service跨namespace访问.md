# Kubernetes-service跨namespace访问

参考文章

1. [kubernetes中跨namespace访问服务](https://blog.csdn.net/jettery/article/details/79226801)

2. [kubernetes: Service located in another namespace](https://blog.csdn.net/jettery/article/details/79226801)

参考文章1, 2中都提到了`ExternalName`这个`service`类型...不过好像没那么麻烦.

不同的namespace中直接通过`service名称.namespace名称.svc.cluster.local`就能访问到不同`namespace`下的指定服务.