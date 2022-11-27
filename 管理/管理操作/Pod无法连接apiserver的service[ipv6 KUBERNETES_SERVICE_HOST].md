# Pod无法连接apiserver的service

## 问题描述

IPv4/IPv6双栈环境中, 某pod启动时, 有如下报错, 并不断重启.

```
panic: unable to load configmap based request-header-client-ca-file: Get https://[2409:808e:4980:330::1:1]:443/api/v1/namespaces/kube-system/configmaps/extension-apiserver-authentication: dial tcp [2409:808e:4980:330::1:1]:443: connect: connection refused
```

其中`[2409:808e:4980:330::1:1]`是apiserver的 service IP 地址(即`default`空间下, 名为`kubernetes`的`Service`资源IP地址)

这个service同时指向了apiserver的IPv4/IPv6地址, 但可能优先解析到了IPv4地址, 而apiserver只监听了IPv6地址, 导致连接失败(或者反过来, service指向了IPv6地址, 而apiserver只监听了IPv4地址).

总之就是, 不是所有的后端apiserver地址, 都能正常访问.

```yaml
apiVersion: v1
kind: Endpoints
metadata:
  name: kubernetes
  namespace: default
subsets:
- addresses:
  - ip: 172.22.248.182
  - ip: 172.22.248.184
  - ip: 2409:808e:4980:730::183
  ports:
  - name: https
    port: 6443
    protocol: TCP
```

## 解决方法

我们知道, Pod内部是可以通过`InClusterConfig()`方法, 得到apiserver的地址的, 这里得到的就是`kubernetes`的`Service`资源IP地址, 我们还可以通过环境变量, 将这个地址覆写, 如下.

```yaml
    env:
    - name: KUBERNETES_SERVICE_HOST
        value: 2409:808e:4980:730::183
    - name: KUBERNETES_SERVICE_PORT
        value: "6443"
```

这样, 容器内的程序在通过`InClusterConfig()`方法, 获取apiserver的地址时, 就是从环境变量里取了.
