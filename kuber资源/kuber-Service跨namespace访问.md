# kuber-通过域名进行跨namespace访问

参考文章

1. [kubernetes中跨namespace访问服务](https://blog.csdn.net/jettery/article/details/79226801)

2. [kubernetes: Service located in another namespace](https://blog.csdn.net/jettery/article/details/79226801)

参考文章1, 2中都提到了`ExternalName`这个`service`类型...不过好像没那么麻烦.

`Pod`内部无法直接访问不同`ns`

不同的`ns`中在Pod直接通过`service名称.ns名称.svc.cluster.local`就能访问到不同`ns`下的指定服务.

同时可以使用`pod名称.service名称.ns名称.svc.cluster.local`来访问相同`ns`下指定的pod.
