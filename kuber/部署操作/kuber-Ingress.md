# kuber-Ingress

参考文章

1. [通俗理解Kubernetes中Service、Ingress与Ingress Controller的作用与关系](https://cloud.tencent.com/developer/article/1326535)
  - 对三者的关系介绍得很清楚
  - ingress暴露给公网访问的3种实践
2. [NGINX Ingress Controller](https://kubernetes.github.io/ingress-nginx/)
  - kubernetes官方组织下只有nginx与gce两种ingress, 这是nginx的组件文档
  - 按照部署文档中的说明, 每个ingress应该是定义单个域名的path规则.
  - [kubectl资源文件](kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/master/deploy/static/mandatory.yaml)

> The default configuration watches Ingress object from all the namespaces. To change this behavior use the flag --watch-namespace to limit the scope to a particular namespace.
> 
> 默认的nginx-ingress-controller可监听所有namespace的ingress对象, 如果需要ta只监听指定namespace, 使用`--watch-namespace`选项.

首先安装nginx-ingress-controller, 见参考文章2.

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: mypod
  labels:
    app: mypod
spec:
  containers:
  - name: nginx-pod
    image: centos:7
    command: ["tail", "-f", "/etc/profile"]
```

注意: 需要进入容器手动安装nginx并启动.

```yaml
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

```yaml
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