yaml文件格式参考

1. [yaml基础语法](http://hustlei.tk/2014/08/yaml-basic-syntax.html)

pod文件参考

1. [官网pod创建](https://kubernetes.io/docs/tasks/configure-pod-container/configure-pod-initialization/)

2. [使用YAML创建一个 Kubernetes Depolyment](https://sanwen.net/a/ytdouqo.html)

```yml
apiVersion: v1
kind: Pod
metadata:
    name: webapp
spec:
    containers:
        - name: nginx
          image: daocloud.io/nginx
```

创建pod

```
[root@localhost pods]# kubectl -s http://172.32.100.90:8080 create -f ./webapp.yml 
pod "webapp" created
[root@localhost pods]# kubectl -s http://172.32.100.90:8080 get pods
NAME      READY     STATUS              RESTARTS   AGE
webapp    0/1       ContainerCreating   0          4m
[root@localhost pods]# kubectl -s http://172.32.100.90:8080 delete pods webapp
pod "webapp" deleted
```

md创建后一直卡在`ContainerCreating`状态, 使用`describe`查看原因

查看pods相关信息

```
kubectl -s http://172.32.100.90:8080 describe pods webapp
```

```
Events:
  FirstSeen	LastSeen	Count	From			SubObjectPath	Type		Reason			Message
  ---------	--------	-----	----			-------------	--------	------			-------
  22s		22s		1	default-scheduler			Normal		Scheduled		Successfully assigned webapp to 172.32.100.70
  22s		22s		1	kubelet, 172.32.100.70			Warning		MissingClusterDNS	kubelet does not have ClusterDNS IP configured and cannot create Pod using "ClusterFirst" policy. Falling back to DNSDefault policy.
```

------

获取service信息

```yml
# Copyright 2016 The Kubernetes Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# __MACHINE_GENERATED_WARNING__

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
        - name: dns
          port: 53
          protocol: UDP
        - name: dns-tcp
          port: 53
          protocol: TCP

```

```
$ kubectl -s http://172.32.100.90:8080 get service
NAME         CLUSTER-IP   EXTERNAL-IP   PORT(S)   AGE
kubernetes   10.10.10.1   <none>        443/TCP   5d
```