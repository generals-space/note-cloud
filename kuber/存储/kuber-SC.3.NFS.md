# kuber-SC.3.NFS

参考文章

1. [external-storage/nfs-client](https://github.com/kubernetes-incubator/external-storage/tree/master/nfs-client)

`externl-storage`项目中有`nfs-client`的子工程, 可以用来当作sc的实际提供者.

在部署之前, 需要先准备好NFS服务端, 然后修改`02.deploy.yaml`文件中的`NFS_SERVER`和`NFS_PATH`. 

前者是NFS服务端的IP地址, 后者则是此NFS服务提供的目录的绝对路径.

该插件的工作原理就是, 在NFS服务端创建一个目录作为根目录, 即`NFS_PATH`所指定的路径. 然后每个指定storage class为此`provisioner`的`PV/PVC`, 在绑定上一个Pod后, 都将在此根目录下创建一个子目录.

