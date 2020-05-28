# CSI插件使用

参考文章

1. [kubernetes-sigs/alibaba-cloud-csi-driver](https://github.com/kubernetes-sigs/alibaba-cloud-csi-driver)
    - 支持的阿里云存储：云盘(Disk)、NAS、CPFS、OSS、LVM
    - 安装上这类插件后, 可以通过在创建PVC时指定StorageClass从而自动创建PV对象, 无需手动操作, 极大减少繁琐的操作.
2. [阿里云Kubernetes CSI实践 - 部署详解](https://yq.aliyun.com/articles/708649)
    - 部署云盘类型的CSI插件
3. [Kubernetes NFS-Client Provisioner](https://github.com/kubernetes-incubator/external-storage/tree/master/nfs-client)

如果了解了kubernetes中`StorageClass`的概念, 使用`CSI`就比较容易了. 

我们知道, 在`SC`的yaml部署文件中可以声明一种`PVC`, 所有指定同一种SC的PVC对象, 都可以按照相同方式从某个地方自动申请, 创建并挂载硬盘资源.

但是`SC`需要一种名为`Provisoner`的程序支持. 比如, 我想定义一个可以自动创建NFS目录类型PV的SC资源, 但起码得有一个装有`nfs-client`的客户端替我去NFS Server创建目录然后再挂载吧? 执行这些操作的组件被称为`Provisioner`. 见参考文章3中的工程, 其本质就是`Provisioner`向指定的NFS Server申请一个目录, 之后创建PVC的请求(指定了SC)都会在这个目录自动创建子目录, 然后创建PV资源绑定这个子目录, 然后被Pod挂载.

阿里云目前拥有的存储方式包括: 云盘(本地硬盘), NAS(NFS服务), OSS(对象存储)等, 他们对每种存储方式都提供了相应的`Provisioner`插件.

以云盘为例, 在当前文件所在目录`CSI插件使用`中提供了部署所用的yaml配置. 

部署环境是阿里云的标准托管kuber集群, 版本为: `1.14.8-aliyun.1`. 

依次创建后可以得到创建好的PVC和PV对象

```console
$ k get pvc
NAME       STATUS   VOLUME                                         CAPACITY   ACCESS MODES   STORAGECLASS                   AGE
disk-pvc   Bound    pv-disk-87a1abb7-1b0c-11ea-8adf-7eab48cdc50c   25Gi       RWO            alicloud-disk-ssd-hangzhou-g   27m
$ k get pv
NAME                                           CAPACITY   ACCESS MODES   RECLAIM POLICY   STATUS   CLAIM              STORAGECLASS                   REASON   AGE
pv-disk-87a1abb7-1b0c-11ea-8adf-7eab48cdc50c   25Gi       RWO            Retain           Bound    default/disk-pvc   alicloud-disk-ssd-hangzhou-g            27m
```

pv名称的生成由插件定义, 目前不清楚具体的格式, 可能可以自行定义.

在云盘控制台, 可以看到与该PV对象对应的云盘资源.

![](https://gitee.com/generals-space/gitimg/raw/master/466D521FF652FDF17189A54843C1B779.png)

------

另外, 云盘也有分类: 高效云盘, ssd云盘, essd云盘, 下面是provisioner组件自行提供的(`alicloud-disk-ssd-hangzhou-g`是我们刚才手动创建的), 其中`alicloud-disk-available`会通过`高效云盘`、`ssd云盘`、`普通云盘`的顺序依次尝试创建云盘.

```console
$ k get sc
NAME                           PROVISIONER                       AGE
alicloud-disk-available        alicloud/disk                     171m
alicloud-disk-efficiency       alicloud/disk                     171m
alicloud-disk-essd             alicloud/disk                     171m
alicloud-disk-ssd              alicloud/disk                     171m
alicloud-disk-ssd-hangzhou-g   diskplugin.csi.alibabacloud.com   132m
```

如果指定`alicloud-disk-efficiency`创建pvc, 可以得到如下结果.

```console
$ k get pvc
NAME                  STATUS   VOLUME                                         CAPACITY   ACCESS MODES   STORAGECLASS                   AGE
disk-pvc              Bound    d-bp18shveit3grcs4sht3                         25Gi       RWO            alicloud-disk-efficiency       2m21s
$ k get pv
NAME                           CAPACITY   ACCESS MODES   RECLAIM POLICY   STATUS   CLAIM                STORAGECLASS                   REASON   AGE
d-bp18shveit3grcs4sht3         25Gi       RWO            Delete           Bound    default/disk-pvc     alicloud-disk-efficiency                3m4s
```

注意点: 

1. 云盘资源在创建时有大小限制, 一般最小为20G, 这个区间可以在控制台查尝试创建一下, 看看具体值.
2. 通过CSI插件创建的 `PVC->PV->[高效云盘, SSD云盘, ESSD云盘]`, 不要直接删除PV对象(删除会卡住), 你需要做的只是删除PVC, 然后PV连同云盘资源也一并会被删除. 如果不小心搞错了顺序, 需要到web控制台确认一下, 可能会有漏掉的情况.
