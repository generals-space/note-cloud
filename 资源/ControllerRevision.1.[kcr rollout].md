# kuber-controller revision

kube: 1.16.2

我试了下, 只有`Daemonset`和`StatefulSet`在创建的时候会自动创建相应的`ControllerRevision`对象, `Deployment`都没有.

另外, 每次对`StatefulSet`对象进行更新操作时, 就会自动创建一个新的`ControllerRevision`对象.

> 以下将`ControllerRevision`简称为`kcr`

```console
$ k get $kcr
NAME                               CONTROLLER                             REVISION   AGE
kube-flannel-ds-amd64-67f65bfbc7   daemonset.apps/kube-flannel-ds-amd64   1          16m
kube-proxy-544458d8d6              daemonset.apps/kube-proxy              1          27m
```

一个`ControllerRevision`对象的内容如下

```yaml
apiVersion: apps/v1
data:
  spec:
    template:
      $patch: replace
      metadata:
        creationTimestamp: null
        labels:
          k8s-app: kube-proxy
      spec:
        ## 省略
kind: ControllerRevision
metadata:
  labels:
    controller-revision-hash: 544458d8d6
    k8s-app: kube-proxy
  name: kube-proxy-544458d8d6
  namespace: kube-system
  ownerReferences:
  - apiVersion: apps/v1
    blockOwnerDeletion: true
    controller: true
    kind: DaemonSet
    name: kube-proxy
    uid: f64e3755-fbd8-41fc-a153-a33e5a98c282
revision: 1
```

`data`块的内容是该`kcr`所属的资源的`spec`部分的内容, 上面我做了一些删减.

> 最新的`kcr`对象内容一般等于当前最新`Daemonset`/`StatefulSet` 对象的内容.
