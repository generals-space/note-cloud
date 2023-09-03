# NFS StorageClass无法挂载nfs-client-provisioner报错unexpected error getting claim reference selfLink was empty, can't make reference

参考文章

1. [问题记录：K8s1.20版本上安装NFS-StorageClass，报错：unexpected error getting claim reference: selfLink was empty.](https://blog.51cto.com/dongweizhen/6316400)
2. [Using Kubernetes v1.20.0, getting "unexpected error getting claim reference: selfLink was empty, can't make reference"](https://github.com/kubernetes-sigs/nfs-subdir-external-provisioner/issues/25)
3. [Kubernetes nfs provider selfLink was empty](https://stackoverflow.com/questions/65376314/kubernetes-nfs-provider-selflink-was-empty)
4. [k8s 1.22使用nfs存储类报错：unexpected error getting claim reference: selfLink was empty, can‘t make referenc](https://blog.51cto.com/zhangxueliang/5647274)
    - 提到了参考文章2

## 问题描述

kube: v1.23.4

按照步骤部署最新的 nfs-client-provisioner 组件, 创建测试 pod 时, pvc 一直处理`Pending`状态, pv 创建不起来.

查看 nfs-client-provisioner 日志, 有如下输出

```log
I1220 22:36:01.615073  1 controller.go:987] provision "default/nfs-client-provisioner" class "nfs-provisioner": started
E1220 22:36:01.618195  1 controller.go:1004] provision "default/nfs-client-provisioner" class "nfs-provisioner": unexpected error getting claim reference: selfLink was empty, can't make reference
```
