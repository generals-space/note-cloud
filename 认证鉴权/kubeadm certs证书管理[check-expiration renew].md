参考文章

1. [Certificate Management with kubeadm](https://kubernetes.io/docs/tasks/administer-cluster/kubeadm/kubeadm-certs/#manual-certificate-renewal)
2. [grpc、https、oauth2等认证专栏实战7：使用cfssl来制作证书介绍](https://blog.csdn.net/u011582922/article/details/126339600)

```
$ kubeadm certs check-expiration
CERTIFICATE                EXPIRES                  RESIDUAL TIME   CERTIFICATE AUTHORITY   EXTERNALLY MANAGED
admin.conf                 Mar 09, 2024 06:28 UTC   364d            ca                      no      
apiserver                  Mar 09, 2024 06:28 UTC   364d            ca                      no      
apiserver-etcd-client      Mar 09, 2024 06:28 UTC   364d            etcd-ca                 no      
apiserver-kubelet-client   Mar 09, 2024 06:28 UTC   364d            ca                      no      
controller-manager.conf    Mar 09, 2024 06:28 UTC   364d            ca                      no      
front-proxy-client         Mar 09, 2024 06:28 UTC   364d            front-proxy-ca          no      
scheduler.conf             Mar 09, 2024 06:28 UTC   364d            ca                      no      

CERTIFICATE AUTHORITY   EXPIRES                  RESIDUAL TIME   EXTERNALLY MANAGED
ca                      Feb 20, 2033 12:52 UTC   9y              no      
etcd-ca                 Feb 20, 2033 12:52 UTC   9y              no      
front-proxy-ca          Feb 20, 2033 12:52 UTC   9y              no  
```

上面列出的crt证书, 及证书的签发机构ca, 对应的都是`/etc/kubernetes/pki`目录下的内容.

```console
$ ll /etc/kubernetes/pki | grep crt
-rw-r--r-- 1 root root 1155 Mar 10 06:28 apiserver-etcd-client.crt
-rw-r--r-- 1 root root 1164 Mar 10 06:28 apiserver-kubelet-client.crt
-rw-r--r-- 1 root root 1289 Mar 10 06:28 apiserver.crt
-rw-r--r-- 1 root root 1099 Feb 23 12:52 ca.crt
-rw-r--r-- 1 root root 1115 Feb 23 12:52 front-proxy-ca.crt
-rw-r--r-- 1 root root 1119 Mar 10 06:28 front-proxy-client.crt
```

使用openssl查看某一证书文件的信息, 是可以对应起来的, 以apiserver.crt文件为例.

```console
$ openssl x509 -in apiserver.crt -noout -text
Certificate:
    Data:
        Issuer: CN = kubernetes
        Validity
            Not Before: Feb 23 12:52:04 2023 GMT
            Not After : Mar  9 06:28:34 2024 GMT
        Subject: CN = kube-apiserver
```

apiserver.crt 由 ca.crt 签发, ta的 Issuer 即为 ca.crt 的 Subject 值.

```console
$ openssl x509 -in ca.crt -noout -text 
Certificate:
    Data:
        Issuer: CN = kubernetes
        Validity
            Not Before: Feb 23 12:52:04 2023 GMT
            Not After : Feb 20 12:52:04 2033 GMT
        Subject: CN = kubernetes
```

## renew 更新证书

```console
$ kubeadm certs renew apiserver 
certificate for serving the Kubernetes API renewed

$ kubeadm certs check-expiration
CERTIFICATE                EXPIRES                  RESIDUAL TIME   CERTIFICATE AUTHORITY   EXTERNALLY MANAGED
admin.conf                 Mar 09, 2024 06:28 UTC   364d            ca                      no      
apiserver                  Mar 09, 2024 06:57 UTC   364d            ca                      no      ## 这个过期时间就与其他不同了.
apiserver-etcd-client      Mar 09, 2024 06:28 UTC   364d            etcd-ca                 no      
apiserver-kubelet-client   Mar 09, 2024 06:28 UTC   364d            ca                      no      
controller-manager.conf    Mar 09, 2024 06:28 UTC   364d            ca                      no      
etcd-healthcheck-client    Mar 09, 2024 06:28 UTC   364d            etcd-ca                 no      
etcd-peer                  Mar 09, 2024 06:28 UTC   364d            etcd-ca                 no      
etcd-server                Mar 09, 2024 06:28 UTC   364d            etcd-ca                 no      
front-proxy-client         Mar 09, 2024 06:28 UTC   364d            front-proxy-ca          no      
scheduler.conf             Mar 09, 2024 06:28 UTC   364d            ca                      no      

CERTIFICATE AUTHORITY   EXPIRES                  RESIDUAL TIME   EXTERNALLY MANAGED
ca                      Feb 20, 2033 12:52 UTC   9y              no      
etcd-ca                 Feb 20, 2033 12:52 UTC   9y              no      
front-proxy-ca          Feb 20, 2033 12:52 UTC   9y              no     
```

此命令会修改`apiserver`的文件内容.

...相关组件是否需要重启???
