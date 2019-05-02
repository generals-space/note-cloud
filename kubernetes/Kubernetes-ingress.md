# Kubernetes-ingress

参考文章

1. [Ingress 支持](https://help.aliyun.com/document_detail/86533.html)

2. [路由配置说明](https://help.aliyun.com/document_detail/86535.html)

3. [部署高可靠 Ingress Controller](https://help.aliyun.com/document_detail/86750.html)

写这篇文章之前, 我还不清楚ingress是什么东西, 只知道service资源用来暴露应用端口, 而ingress似乎可以实现类似于nginx根据域名, path定义不同的转发规则.

首先按照参考文章1和2阿里云提供的文档, 把ingress实例先创建起来, 再回来查看比较官方的文档.

> 写在前面: 失败了, ingress创建完成后访问地址返回404, 提交工单给阿里云, 回复说集群版本(1.9.7)太低, `nginx-ingress-controller`组件不支持. 参考文章1和2根本没提到过`ingress-controller`, 当支持人员询问这个组件的版本时我蒙了一下. 好在文档下方有参考文章3的链接, `nginx-ingress-controller`是一个存在于`kube-system`命名空间下的pod对象, 使用`describe`命令可以查看ta的详细信息.

```yml
apiVersion: v1
kind: Pod
metadata:
  name: mypod
spec:
  containers:
  - name: nginx-pod
    image: centos:7
    command: ["tail", "-f", "/etc/profile"]
```

```yml
apiVersion: v1
kind: Service
metadata:
  name: myservice
spec:
  selector:
    app: mypod
  ports:
  - name: default-port
    protocol: TCP
    port: 80
    targetPort: 80
  type: ClusterIP
```

```yml
apiVersion: extensions/v1beta1
kind: Ingress
metadata:
  name: myingress
  annotations:
    nginx.ingress.kubernetes.io/rewrite-target: /
spec:
  rules:
  - host: foo.example.com
    http:
      paths:
      - path: /svc
        backend:
          serviceName: myservice
          servicePort: 80

```

ok, 创建完成后, 查看ingress资源, 我们发现ADDRESS字段已经有一个值了.

```
$ k get ingress
NAME        HOSTS             ADDRESS         PORTS     AGE
myingress   foo.example.com   106.14.49.105   80        22m
```

查看阿里云后台

![](https://gitee.com/generals-space/gitimg/raw/master/e161c34e6ad6adc8c58344bec15cb515.png)

...访问404了