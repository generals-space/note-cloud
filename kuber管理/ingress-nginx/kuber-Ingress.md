# kuber-Ingress

参考文章

1. [通俗理解Kubernetes中Service、Ingress与Ingress Controller的作用与关系](https://cloud.tencent.com/developer/article/1326535)
  - 对三者的关系介绍得很清楚
  - ingress暴露给公网访问的3种实践
2. [NGINX Ingress Controller](https://kubernetes.github.io/ingress-nginx/)
  - kubernetes官方组织下只有nginx与gce两种ingress, 这是nginx的组件文档
  - 按照部署文档中的说明, 每个ingress应该是定义单个域名的path规则.
  - [kubectl资源文件](kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/master/deploy/static/mandatory.yaml)
3. [官方文档 ingress-nginx Annotations](https://github.com/kubernetes/ingress-nginx/blob/master/docs/user-guide/nginx-configuration/annotations.md)
    - ingress 所有可用注解
4. [How to proxy_pass with nginx-ingress?](https://stackoverflow.com/questions/57806618/how-to-proxy-pass-with-nginx-ingress)
    - ingress 里没有`proxy_pass`, 只能用`write-target`替代.

> The default configuration watches Ingress object from all the namespaces. To change this behavior use the flag --watch-namespace to limit the scope to a particular namespace.
> 
> 默认的nginx-ingress-controller可监听所有namespace的ingress对象, 如果需要ta只监听指定namespace, 使用`--watch-namespace`选项.

首先安装nginx-ingress-controller, 见参考文章2. ingress-controller将作为一个pod运行在**某个worker节点**上.

然后依次创建pod, service, ingress对象.

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
    image: nginx
```

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
    ## 这一句注解可以将发送到ingress controller的/svc的请求, 转发给后端Service时路径转为/
    ## 类似nginx在proxy_pass路径末尾的斜线.
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

------

补充, 更进一步.

上面的`ingress`配置, 对于所有的`/svc`, 以及`/svc/xxx`请求, 最终都只请求`myservice/`根路径, 如果`myservice`提供了多个路径, 就显得捉襟见肘了, 所以需要改成如下配置.

```yaml
apiVersion: extensions/v1beta1
kind: Ingress
metadata:
  name: myingress
  annotations:
    nginx.ingress.kubernetes.io/rewrite-target: /$2
spec:
  rules:
  - host: foo.example.com
    http:
      paths:
      - path: /svc(/|$)(.*)
        backend:
          serviceName: myservice
          servicePort: 80
```

在原生 nginx 中, proxy_pass 与 rewrite 的含义是完全不同的, 前者表示向后端`upstream`池中进行转发, 而后者则是向浏览器返回 301/302 的重定向响应, 然后浏览器再次发起请求.

很明显, ingress 是把`rewrite-target`当成`proxy_pass`用了.
