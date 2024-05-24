# readOnlyRootFilesystem以只读方式加载容器的根文件系统

参考文章

1. [Configure a Security Context for a Pod or Container](https://kubernetes.io/docs/tasks/configure-pod-container/security-context/)
    - 官方文档
2. [为 Pod 或容器配置安全上下文](https://kubernetes.io/zh-cn/docs/tasks/configure-pod-container/security-context/)
    - 官方文档
    - 参考文章1的中文版
    - `readOnlyRootFilesystem: true`, 以只读方式加载容器的根文件系统(默认为`false`)

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  namespace: default
  name: mydeploy
  labels:
    app: mydeploy
spec:
  replicas: 1
  selector:
    matchLabels:
      app: mydeploy
  template:
    metadata:
      labels:
        app: mydeploy
    spec:
      containers:
      - image: busybox:1.32.0
        command: ['sh', '-c', 'tail -f /dev/null']
        imagePullPolicy: IfNotPresent
        name: busybox
        securityContext:
          readOnlyRootFilesystem: true
        volumeMounts:
        - mountPath: /tmp
          name: temp-vol
      volumes:
      - name: temp-vol
        emptyDir: {}
```

配置`readOnlyRootFilesystem: true`后, 容器内部无法对文件系统做任何修改, 但是无法影响通过`volumeMounts`挂载的目录.

```log
$ k -n default exec mydeploy-5b7ff8d464-j2nrj -- touch /abc.txt
touch: /abc.txt: Read-only file system
command terminated with exit code 1

$ k -n default exec mydeploy-5b7ff8d464-j2nrj -- touch /var/abc.txt
touch: /var/abc.txt: Read-only file system
command terminated with exit code 1

$ k -n default exec mydeploy-5b7ff8d464-j2nrj -- touch /etc/abc.txt
touch: /etc/abc.txt: Read-only file system
command terminated with exit code 1

$ k -n default exec mydeploy-5b7ff8d464-j2nrj -- touch /tmp/abc.txt

$ k -n default exec mydeploy-5b7ff8d464-j2nrj -- ls /tmp
abc.txt
```
