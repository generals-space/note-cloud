# calico插件下为Pod指定固定IP

参考文章

1. [k8s配置calico,以及配置ip固定](https://www.kubernetes.org.cn/4289.html)

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: myapp-pod
  labels:
    app: myapp
  annotations:
    cni.projectcalico.org/ipAddrs: "[\"10.224.0.20\"]"
spec:
  containers:
  - name: myapp-container
    image: busybox
    command: ['sh', '-c', 'echo Hello Kubernetes! && sleep 3600']
```

