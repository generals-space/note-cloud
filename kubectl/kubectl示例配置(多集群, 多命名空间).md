# kubectl示例配置(多集群, 多命名空间)

<!--
<!tags!>: <!kubernetes!>
<!keys!>: Pyak5endi8nCjjr[
-->
```yml
apiVersion: v1
clusters:
  - cluster:
      certificate-authority-data: // 这里是base64加密过的字符串
      server: https://106.14.49.2:6443
  name: kubernetes
  - cluster:
      certificate-authority: /Users/general/.minikube/ca.crt
      server: https://192.168.64.32:8443
  name: minikube

contexts:
  - context:
      cluster: kubernetes
      namespace: lora-app
      user: kubernetes-admin
    name: aliyun
  - context:
      cluster: minikube
      user: minikube
    name: minikube

current-context: aliyun
kind: Config
preferences: {}
users:
  - name: kubernetes-admin
    user:
      client-certificate-data: // 这里是base64加密过的字符串
      client-key-data: // 这里是base64加密过的字符串
  - name: minikube
    user:
      client-certificate: /Users/general/.minikube/client.crt // 也可以写文件路径
      client-key: /Users/general/.minikube/client.key // 也可以写文件路径
```

上述配置文件包含了两个集群(Cluster): `kubernetes(运行在阿里云)`和`minikube(运行在本地)`.

每个集群又包含n个工作环境(Context): `kubernetes`集群的是`aliyun`, `minikube`集群的是同名的`minikube`.

每个工作环境又可以包含n个命名空间(Namespace), `aliyun`下的叫`lora-app`(这是我手动建的), `minikube`下的没写, 默认为`default`.

按照概念大小来看, Cluster > Context > Namespace.

- Cluster是物理集群, 每一项都需要一个`apiserver`的http接口表示

- Context表示工作环境, 作用上相当于`develope`, `test`, `product`这种.

- Namespace则可以用于划分业务, 不同业务可以隔离, 不会通过service暴露的接口直接通信, 但仍有办法通过dns解析得到其他命名空间的服务路径.

注意: namespace并不是属于context的, 它们并不是包含关系, 而是绑定关系, 一个namespace可以绑定在**context1**上, 也可以绑定到**context2**上.