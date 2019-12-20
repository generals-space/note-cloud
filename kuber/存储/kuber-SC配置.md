# kuber-SC本地卷LocalPV(一)

```yaml
# Only create this for K8s 1.9+
apiVersion: storage.k8s.io/v1
kind: StorageClass
metadata:
  name: fast-disks
provisioner: kubernetes.io/no-provisioner
volumeBindingMode: WaitForFirstConsumer
# Supported policies: Delete, Retain
reclaimPolicy: Delete
```

## volumeBindingMode

`WaitForFirstConsumer`: 单纯创建pvc时不会自动创建pv, 也不会创建实际目录, 此时pvc的状态保持为`Pending`. 只有当pod引用此pvc时才会创建, 同时pvc的状态会变为`Bound`. 不过pod被删除后, 遗留的pvc仍然为`Bound`状态.

