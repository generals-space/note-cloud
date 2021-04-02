# calico-为Pod指定静态IP [ipAddrs]

参考文章

1. [k8s配置calico,以及配置ip固定](https://www.kubernetes.org.cn/4289.html)
2. [calico官方文档 - Use a specific IP address with a pod](https://docs.projectcalico.org/networking/use-specific-ip)

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: centos7-pod
  labels:
    app: centos7
  annotations:
    cni.projectcalico.org/ipAddrs: '["10.224.0.20"]'
spec:
  containers:
  - name: centos7
    image: registry.cn-hangzhou.aliyuncs.com/generals-space/centos7
    command: ['tail', '-f', '/etc/os-release']
```

`ipAddrs`可以是合法的, PodCIDR 中的任意 IP, 且**可以是当前集群中未被划分的网段中的 IP**. 

```
$ k get pod -o wide
NAME            READY   STATUS    RESTARTS   AGE     IP               NODE            NOMINATED NODE   READINESS GATES
centos7-pod     1/1     Running   0          4m17s   172.16.240.100   k8s-worker-01   <none>           <none>
test-ds-26xf9   1/1     Running   0          38h     172.16.151.128   k8s-master-01   <none>           <none>
test-ds-2rr9q   1/1     Running   0          38h     172.16.36.193    k8s-worker-01   <none>           <none>
test-ds-mc7m8   1/1     Running   0          38h     172.16.118.64    k8s-worker-02   <none>           <none>
```

上面的`172.16.240.100`不属于3个节点中任何一个网段, 这种情况下 calico 将会进行随机调度, 其他节点将添加到此 IP 的黑洞路由`blackhole 172.16.240.64/26 proto bird`.

另外, `ipAddrs`特性貌似只能用于 Pod 资源, 无法实现 [cni-terway](https://github.com/generals-space/cni-terway) 的`Deployment/DaemonSet` IP 池的功能.
