# kuber-不通过PV与PVC直接挂载NFS目录

参考文章

1. [Configuring NFS Storage for Kubernetes](https://docs.docker.com/ee/ucp/kubernetes/storage/use-nfs-volumes/)

```yaml
kind: Pod
apiVersion: v1
metadata:
  name: nfs-in-a-pod
spec:
  containers:
    - name: app
      image: alpine
      volumeMounts:
        - name: nfs-volume
          mountPath: /var/nfs # Please change the destination you like the share to be mounted too
      command: ["/bin/sh"]
      args: ["-c", "sleep 500000"]
  volumes:
    - name: nfs-volume
      nfs:
        server: nfs.example.com # Please change this to your NFS server
        path: /share1 # Please change this to the relevant share
```

注意: 运行Pod的worker节点需要安装`nfs-utils`工具, 否则Pod会无法启动, 一直卡住.
