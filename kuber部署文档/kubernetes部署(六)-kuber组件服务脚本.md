---
title: kubernetes部署(六)-kuber组件服务脚本
tags: [kubernetes]
categories: general
---

<!--

# kubernetes部署(六)-kuber组件服务脚本

<!tags!>: <!kubernetes!>

<!keys!>: jfQ;ks8fqJm3dcdw

-->

## 1. apiserver

`--etcd-servers`是指定的参数, 它表示了etcd集群/单点的访问路径, 这里没有配置https, 所以没有指定证书路径.

`--insecure-bind-address`与`--insecure-port=8080`指定了apiserver服务监听的地址跟端口, 同样是非https, 否则将会是6443.

`--service-cluster-ip-range`指定了在kuber集群中, service对象分配的ip的网段. (service与pods是同级概念, 但作用不同, 详情请参考[Kubernetes 中的PodIP、ClusterIP 和外部IP](http://blog.csdn.net/liukuan73/article/details/54773579))

`--kubelet-https`默认为true, 提供https接口, 在摸索阶段还是设置为false吧.

```bash
## /usr/lib/systemd/system/apiserver.service 
[Unit]
Description=Kubernetes API Server
Documentation=https://github.com/GoogleCloudPlatform/kubernetes
[Service]
ExecStart=/usr/local/kubernetes/bin/kube-apiserver  \
    --etcd-servers=http://172.32.100.71:2379 \
    --insecure-bind-address=172.32.100.71 \
    --insecure-port=8080 \
    --service-cluster-ip-range=10.10.10.0/24 \
    --kubelet-https=false \
    --allow-privileged=false \
    --logtostderr=true \
    --v=4
Restart=on-failure
[Install]
WantedBy=multi-user.target
```

## 2. scheduler

唯一的要求也就是`--master`选项了, 指定apiserver的访问地址.

```bash
## /usr/lib/systemd/system/scheduler.service 
[Unit]
Description=Kubernetes Scheduler
Documentation=https://github.com/GoogleCloudPlatform/kubernetes
[Service]
ExecStart=/usr/local/kubernetes/bin/kube-scheduler \
    --master=172.32.100.71:8080 \
    --logtostderr=true \
    --v=4
Restart=on-failure
[Install]
WantedBy=multi-user.target
```

## 3. controller-manager

同样, 只需要指定apiserver的访问地址.

```bash
## /usr/lib/systemd/system/controller-manager.service 
[Unit]
Description=Kubernetes Controller Manager
Documentation=https://github.com/GoogleCloudPlatform/kubernetes
[Service]
ExecStart=/usr/local/kubernetes/bin/kube-controller-manager \
    --master=172.32.100.71:8080 \
    --logtostderr=true \
    --v=4
Restart=on-failure
[Install]
WantedBy=multi-user.target
```

## 4. proxy

只要指定apiserver...突然发现kubernetes的集群配置好简单啊怎么办...

```bash
## /usr/lib/systemd/system/proxy.service 
[Unit]
Description=Kubernetes Proxy
# the proxy crashes if etcd isn't reachable.
# https://github.com/GoogleCloudPlatform/kubernetes/issues/1206
After=network.target
[Service]
ExecStart=/usr/local/kubernetes/bin/kube-proxy \
    --master=172.32.100.71:8080 \
    --logtostderr=true \
    --v=4
Restart=on-failure
[Install]
WantedBy=multi-user.target
```

## 5. kubelet

同样, 配置apiserver地址.

然后

`--cgroup-driver`设置为systemd, 默认应该是`cgroupfs`, 未设置这个选项时我这边出了点总是, kubelet无法启动...忘了截图, 算了.

`--cluster-dns`与`--cluster-domain=cluster.local`选项是为了之后安装dns插件准备的, 因为在apiserver中的`--service-cluster-ip-range`选项中设置了`10.10.10.0/24`网段, 而kuber会在这个网段把它自己启动成一个服务, 并且占用的是`10.10.10.1`, 所以给dns服务预留一个`10.10.10.2`.

```bash
## /usr/lib/systemd/system/kubelet.service 
[Unit]
Description=Kubernetes Kubelet
After=docker.socket cadvisor.service
## Requires=docker.socket
[Service]
ExecStart=/usr/local/kubernetes/bin/kubelet \
    --api-servers=172.32.100.71:8080 \
    --cgroup-driver=systemd \
    --cluster-domain=cluster.local \
    --cluster-dns=10.10.10.2 \
    --allow-privileged=false \
    --logtostderr=true \
    --v=4
Restart=on-failure
[Install]
WantedBy=multi-user.target
```
