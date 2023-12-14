# kuber-通过域名进行跨namespace访问

参考文章

1. [kubernetes中跨namespace访问服务](https://blog.csdn.net/jettery/article/details/79226801)
2. [kubernetes: Service located in another namespace](https://blog.csdn.net/jettery/article/details/79226801)
    - service名称.ns名称.svc.cluster.local
3. [k8s教程（service篇）-pod的dns域名](https://blog.csdn.net/qq_20042935/article/details/128674001)
    - pod-ip.ns名称.pod.cluster.local

参考文章1, 2中都提到了`ExternalName`这个`service`类型...不过好像没那么麻烦.

不同`ns`中的 pod, 直接通过`service名称.ns名称.svc.cluster.local`就能访问到不同`ns`下的指定服务, 要求被访者拥有 service 资源.

如果是`sts`类型生成的 pod, 还可以通过`pod名称.service名称.ns名称.svc.cluster.local`来访问单一 pod.

------

但是如果是纯 pod, 或是 deployment/daemonset 生成的 pod, 是无法直接通过 pod 名称访问到指定 pod 的(pod 名称后面跟着随意字符串, 没法固定).

如果 pod 由 deployment/daemonset 生成, 也可以用`10-0-95-63.deployment/daemonset名称.default.svc.cluster.local`代替.

对于纯 pod, 则可以使用`10-0-95-63.default.pod.cluster.local`访问, 即"podip(把点号替换成中横线).ns名称.pod.cluster.local".

------

但 podip 是在 pod 创建完成后分配的, 想访问目标 pod 就只能等ta先启动, 这样就需要严格控制启动顺序, 很不方便.

参考文章3中提供了一种方法

```yaml
apiversion: v1
kind: Pod
metadata:
  name: webapp1
  labels:
    app: webapp1
spec:
  hostname: webapp-1 
  subdomain: mysubdomain
  containers:
  - name: webapp1
    image: kubeguide/tomcat-app: v1 
  ports:
  - containerPort: 8080
```

为 pod 设置`hostname`和`subdomain`字段后, 可以通过`webapp-1.mysubdomain.default.svc.cluster.local`访问该 pod.

如果把`subdomain`设置为对应的 service 资源的名称, 把`hostname`设置为 pod 的名称, 就可以与 sts 生成的 pod 用同样的访问方式了.
