# kuber-Ingress开启https

参考文章

1. [Kubernetes 使用 ingress 配置 https 集群(十五)](https://www.cnblogs.com/wzlinux/p/10159366.html)

2. [NGINX Ingress Controller - TLS/HTTPS](https://kubernetes.github.io/ingress-nginx/user-guide/tls/)

与nginx类似, 用ingress开启https也需要在ingress对象的配置文件中指定证书与密钥. 不过不同的是, 在kuber中, 证书与密钥是以`tls`形式的`secret`对象存在的, 需要读入crt和key文件以创建.

不过需要注意的是, tls需要的证书与密钥对貌似需要特定的格式, 我尝试把为nginx生成的证书(在nginx中实践成功)创建为tls对象, 但是客户端访问时出现

```

```

ingress-controller的pod中日志如下.

```
2019/09/03 10:44:17 [error] 42#42: *120 [lua] certificate.lua:79: call(): failed to set DER cert: SSL_use_certificate() failed, context: ssl_certificate_by_lua*, client: 192.168.124.85, server: 0.0.0.0:443
2019/09/03 10:44:17 [crit] 42#42: *118 SSL_do_handshake() failed (SSL: error:1417A179:SSL routines:tls_post_process_client_hello:cert cb error) while SSL handshaking, client: 192.168.124.85, server: 0.0.0.0:443
```

在谷歌上查阅了很久, 之前以为可能是ingress-controller版本的问题, 从1.24.1升级到了1.25.0, 但是没用.

后来怀疑是证书格式的问题, 按照参考文章2中, 也就是官方文档中生成证书的方法, 终于成功了.

```console
openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout server.key -out server.crt -subj "/CN=harbor.generals.space/O=harbor.generals.space"

kubectl create secret tls https-certs -n harbor --key server.key --cert server.crt
```

```yml
apiVersion: extensions/v1beta1
kind: Ingress
metadata:
  name: harbor-ing
  namespace: harbor
spec:
  tls:
  - hosts:
    - harbor.generals.space
    secretName: https-certs
  rules:
  - host: harbor.generals.space
    http:
      paths:
      - path: /
        backend:
          serviceName: portal-svc
          servicePort: 80
```

实验成功.
