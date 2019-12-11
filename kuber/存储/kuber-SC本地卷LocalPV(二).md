# kuber-SC本地卷LocalPV(一)

参考文章

1. [rancher/local-path-provisioner](https://github.com/rancher/local-path-provisioner)

在进行本地测试时, 经常需要对数据做一些持久化, 但又不是太重要的数据, 用hostPath足矣. PV资源与PVC一样可以在yaml文件中定义, 但绑定的Pod必须在PV指向的路径确切存在的时候才能正常启动, 否则就会Crash, 这就要求在创建PV的时候事先手动创建宿主机目录. 用helm测试时(很多helm包都需要做持久化), 一定要创建PVC, 如果频繁操作简直要被烦死.

rancher提供了一个`local pv`的`provisioner`插件, 其本质也是用`hostPath`做的. 在创建PVC的时候, 只要指定rancher提供的sc, 就可以自动创建PV且创建宿主机目录, 然后可以直接被挂载到Pod对象中.

```console
$ k get sc
NAME         PROVISIONER             AGE
local-path   rancher.io/local-path   44m
```

现在可以使用了, `11.pvc.yaml`和`12.deploy.yaml`是测试文件.

pvc和deploy创建成功后, Pod被调度到哪个节点上, provisioner就会在此节点上的`/opt/local-path-provisioner`目录下创建子目录, 同时也会创建PV资源, 目录名称与PV对象的名称相同.

```console
$ k get pod -o wide
NAME                              READY   STATUS    RESTARTS   AGE     IP             NODE            NOMINATED NODE   READINESS GATES
local-pv-deploy-79c9f76b9-fcxd5   1/1     Running   0          20m     10.23.36.247   k8s-worker-01   <none>           <none>
$ k get pvc
NAME        STATUS   VOLUME                                     CAPACITY   ACCESS MODES   STORAGECLASS   AGE
local-pvc   Bound    pvc-90ab9d86-cb83-407d-823c-f7e62b9d5636   1Gi        RWO            local-path     15m
$ k get pv
NAME                                       CAPACITY   ACCESS MODES   RECLAIM POLICY   STATUS   CLAIM               STORAGECLASS   REASON   AGE
pvc-90ab9d86-cb83-407d-823c-f7e62b9d5636   1Gi        RWO            Retain           Bound    default/local-pvc   local-path              15m
```

```console
$ pwd
/opt/local-path-provisioner
$ ls
pvc-90ab9d86-cb83-407d-823c-f7e62b9d5636
```

需要注意的是, 使用说明中给出的yaml部署文件, sc的回收策略`reclaimPolicy`定义为`Delete`, 即绑定的PVC被删除时, 自动创建的PV也会被删除, 这样是非常危险的, 所以我将其修改为`Retain`. 这样, 在pvc被删除后, pv也会对应删除, 但是宿主机上的目录是不会被删除的.

