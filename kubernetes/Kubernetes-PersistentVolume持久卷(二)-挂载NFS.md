# Kubernetes-PersistentVolume持久卷(二)-挂载NFS

参考文章

1. [Kubernetes使用NFS作为共享存储](https://blog.51cto.com/passed/2160149)

2. [KUBERNETES存储之PERSISTENT VOLUMES简介](https://www.cnblogs.com/styshoo/p/6731425.html)

NFS是外部资源, 就是事先拥有一个NFS服务器, 然后kubernetes可以通过创建`PV`资源来使用. 创建PV和PVC的过程和常规的`hostPath`一样.

假设在服务器`192.168.1.100:/mnt`上有一个可用的NFS服务.

### 1. 创建pv资源

```yml
## pv
kind: PersistentVolume
apiVersion: v1
metadata:
  name: nfs-volume
  labels:
    type: nfs
spec:
  ## storageClassName用于绑定pv和pvc, 可以自定义.
  storageClassName: nfs-pv
  capacity:
    storage: 80Gi
  accessModes:
    - ReadWriteMany
  nfs: #############################################只有这里不一样.
    path: "/mnt/nfs"
    server: 192.168.1.100
---

## pvc
kind: PersistentVolumeClaim
apiVersion: v1
metadata:
  name: nfs-claim
spec:
  storageClassName: nfs-pv
  accessModes:
    - ReadWriteMany
  resources:
    requests:
      storage: 80Gi
```

### 2. 使用资源, 常规的pod添加上`volumes`声明即可.

```yml
## fortest
---
apiVersion: v1
kind: Pod
metadata:
  name: fortest
spec:
  volumes:
    - name: fortest-vol
      persistentVolumeClaim:
        claimName: nfs-claim  ## 注意这里.
  containers:
  - name: fortest
    image: generals/golang
    imagePullPolicy: Always
    volumeMounts:
    - name: fortest-vol
      mountPath: "/upload"
      readOnly: false ## 默认为true???
    command: ["tail", "-f", "/etc/profile"]
```

> 注意: deployment文件中`readOnly: false`是必须的, 貌似默认为true, 但是pod的配置默认为false, 可写. 所以deployment的配置要注意一下.