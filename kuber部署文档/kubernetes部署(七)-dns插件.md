---
title: kubernetes部署(七)-dns插件
tags: [kubernetes]
categories: general
---

<!--

# kubernetes部署(七)-dns插件

<!tags!>: <!kubernetes!>

<!keys!>: 3leeU1dijlsCzvk*

-->


参考文章

1. [Kubernetes集群DNS插件安装](http://tonybai.com/2016/10/23/install-dns-addon-for-k8s/)

2. [docker 使用 socks5代理](http://www.jianshu.com/p/fef11e46ebf1)

dns插件以pod与service形式运行在kuber集群中. 这样就不免用到yaml配置文件.

因为我是编译安装, 所以下载了`kubernetes`源码, dns插件的模板文件在`/root/gopath/src/k8s.io/kubernetes/cluster/addons/dns`. 如果没有下载下来, 可以参考github上kubernetes仓库中`kubernetes/cluster/addons/dns`目录.

其中有很多`.base`, `.in`, `.sed`为结尾的文件, 我们这种安装方式, 只用到如下4个配置文件.

```
-rw-r--r--. 1 root root  731 Jul  3 02:50 kubedns-cm.yaml
-rw-r--r--. 1 root root 5327 Jul  3 02:50 kubedns-controller.yaml.base
-rw-r--r--. 1 root root  187 Jul  3 02:50 kubedns-sa.yaml
-rw-r--r--. 1 root root 1037 Jul  3 02:50 kubedns-svc.yaml.base
```

其中`.base`是模板文件, 需要手动修改一些东西, `kubedns-cm.yaml`与`kubedns-sa.yaml`无需做修改, 可以直接创建.

```
$ kubectl -s http://172.32.100.71:8080 create -f kubedns-cm.yaml
$ kubectl -s http://172.32.100.71:8080 create -f kubedns-sa.yaml
```

cm是ConfigMap的缩写, sa是ServiceAccount的缩写, 正是`kubedns-controller.yaml`中`configMap`与`serviceAccountName`字段需要用到的配置...虽然并不知道有什么用.

然后是`kubedns-svc.yaml`文件, 只修改了其中的`clusterIP`字段, 它表示集群中dns服务的ip地址, 也正好是Minion节点上`kubelet`服务中`--cluster-dns`选项的值. 


```yaml
## kubedns-svc.yaml
apiVersion: v1
kind: Service
metadata:
    name: kube-dns
    namespace: kube-system
    labels:
        k8s-app: kube-dns
        kubernetes.io/cluster-service: "true"
        addonmanager.kubernetes.io/mode: Reconcile
        kubernetes.io/name: "KubeDNS"
spec:
    selector:
        k8s-app: kube-dns
    clusterIP: 10.10.10.2
    ports:
        -   
            name: dns
            port: 53
            protocol: UDP
        -   
            name: dns-tcp
            port: 53
            protocol: TCP
```

修改完成后, 按照上面的方法创建.

```
$ kubectl -s http://172.32.100.71:8080 create -f ./kubedns-svc.yaml 
service "kube-dns" created
$ kubectl -s http://172.32.100.71:8080 get services --all-namespaces
NAMESPACE     NAME         CLUSTER-IP   EXTERNAL-IP   PORT(S)         AGE
default       kubernetes   10.10.10.1   <none>        443/TCP         24m
kube-system   kube-dns     10.10.10.2   <none>        53/UDP,53/TCP   40s
```

然后是`kubedns-controller.yaml`. 将其中的`__PILLAR__DNS__DOMAIN__`替换成`10.10.10.2`, 也即kubelet服务`--cluster-domain`选项的值. 这个值其实可以自定义, 保持一致就行, 这里延用网络上其他教程的`cluster.local`.

除了这个, 很重要的一点就是, `kubedns`容器的args选项中`--kube-master-url=http://172.32.100.71:8080`的设置, 因为我们没有使用证书, 所以apiserver的默认端口6443是无法使用的, 而`kubedns`容器默认连接的就是apiserver的6443, 这会导致容器无法启动.

经历了无数次失败后才找到了参数文章1关于这一点的讲解, 非常感谢原作者.

还有, 因为`kubedns-controller`用到了3个镜像, 需要从`gcr.io`下载. 你需要办法下载到本地, 我用的是本地的socks5代理事先将镜像下载到Minion节点上的.

关于docker pull时socks5代理的设置方法, 见参考文章2.(不过里面设置的白名单好像没用, 最好下载完毕后恢复到原来的设置)

```yaml
## kubedns-controller.yaml
apiVersion: extensions/v1beta1
kind: Deployment
metadata:
    name: kube-dns
    namespace: kube-system
    labels:
        k8s-app: kube-dns
        kubernetes.io/cluster-service: "true"
        addonmanager.kubernetes.io/mode: Reconcile
spec:
    strategy:
        rollingUpdate:
            maxSurge: 10%
            maxUnavailable: 0
    selector:
        matchLabels:
            k8s-app: kube-dns
    template:
        metadata:
            labels:
                k8s-app: kube-dns
            annotations:
                    scheduler.alpha.kubernetes.io/critical-pod: ''
        spec:
            tolerations:
                - 
                    key: "CriticalAddonsOnly"
                    operator: "Exists"
            volumes:
                - 
                    name: kube-dns-config
                    configMap:
                        name: kube-dns
                        optional: true
            containers:
                - 
                    name: kubedns
                    image: gcr.io/google_containers/k8s-dns-kube-dns-amd64:1.14.4
                    resources:
                        limits:
                            memory: 170Mi
                        requests:
                            cpu: 100m
                            memory: 70Mi
                    livenessProbe:
                        httpGet:
                            path: /healthcheck/kubedns
                            port: 10054
                            scheme: HTTP
                        initialDelaySeconds: 60
                        timeoutSeconds: 5
                        successThreshold: 1
                        failureThreshold: 5
                    readinessProbe:
                        httpGet:
                            path: /readiness
                            port: 8081
                            scheme: HTTP
                        initialDelaySeconds: 3
                        timeoutSeconds: 5
                    args:
                        - --domain=cluster.local.
                        - --dns-port=10053
                        - --config-dir=/kube-dns-config
                        - --v=2
                        - --kube-master-url=http://172.32.100.71:8080 
                    env:
                        - 
                            name: PROMETHEUS_PORT
                            value: "10055"
                    ports:
                        - 
                            containerPort: 10053
                            name: dns-local
                            protocol: UDP
                        - 
                            containerPort: 10053
                            name: dns-tcp-local
                            protocol: TCP
                        - 
                            containerPort: 10055
                            name: metrics
                            protocol: TCP

                - 
                    name: dnsmasq
                    image: gcr.io/google_containers/k8s-dns-dnsmasq-nanny-amd64:1.14.4
                    livenessProbe:
                        httpGet:
                            path: /healthcheck/dnsmasq
                            port: 10054
                            scheme: HTTP
                        initialDelaySeconds: 60
                        timeoutSeconds: 5
                        successThreshold: 1
                        failureThreshold: 5
                    args:
                        - -v=2
                        - -logtostderr
                        - -configDir=/etc/k8s/dns/dnsmasq-nanny
                        - -restartDnsmasq=true
                        - --
                        - -k
                        - --cache-size=1000
                        - --log-facility=-
                        - --server=/cluster.local/127.0.0.1#10053
                        - --server=/in-addr.arpa/127.0.0.1#10053
                        - --server=/ip6.arpa/127.0.0.1#10053
                    ports:
                        - 
                            containerPort: 53
                            name: dns
                            protocol: UDP
                        - 
                            containerPort: 53
                            name: dns-tcp
                            protocol: TCP
                    resources:
                        requests:
                            cpu: 150m
                            memory: 20Mi
                    volumeMounts:
                        - 
                            name: kube-dns-config
                            mountPath: /etc/k8s/dns/dnsmasq-nanny
                - 
                    name: sidecar
                    image: gcr.io/google_containers/k8s-dns-sidecar-amd64:1.14.4
                    livenessProbe:
                        httpGet:
                            path: /metrics
                            port: 10054
                            scheme: HTTP
                        initialDelaySeconds: 60
                        timeoutSeconds: 5
                        successThreshold: 1
                        failureThreshold: 5
                    args:
                        - --v=2
                        - --logtostderr
                        - --probe=kubedns,127.0.0.1:10053,kubernetes.default.svc.cluster.local,5,A
                        - --probe=dnsmasq,127.0.0.1:53,kubernetes.default.svc.cluster.local,5,A
                    ports:
                        - 
                            containerPort: 10054
                            name: metrics
                            protocol: TCP
                    resources:
                        requests:
                            memory: 20Mi
                            cpu: 10m
            dnsPolicy: Default  # Don't use cluster DNS.
            serviceAccountName: kube-dns
```

部署

```
$ kubectl -s http://172.32.100.71:8080 create -f ./kubedns-controller.yaml 
$ kubectl -s http://172.32.100.71:8080 get deploy --all-namespaces
NAMESPACE     NAME       DESIRED   CURRENT   UP-TO-DATE   AVAILABLE   AGE
kube-system   kube-dns   1         1         1            1           45s
```

