# local-pv部署文件中accessMode不可声明为ReadWriteMany

参考文章

1. [Why the provisioner only support ReadWriteOnce](https://github.com/rancher/local-path-provisioner/issues/70)

实验用的PVC部署文件中的`accessMode`字段只能使用`ReadWriteOnce`, 如果是`ReadWriteMany`会导致PV资源无法创建. 查看`provisioner`的日志会发现如下报错.

```
ERROR: logging before flag.Parse: I0314 14:56:13.629550       1 event.go:221] Event(v1.ObjectReference{Kind:"PersistentVolumeClaim", Namespace:"default", Name:"redis-data-redis-app-0", UID:"d8d33f38-43b7-45b3-b77d-31740e25f4d6", APIVersion:"v1", ResourceVersion:"337966", FieldPath:""}): type: 'Warning' reason: 'ProvisioningFailed' failed to provision volume with StorageClass "local-path": Only support ReadWriteOnce access mode
```
