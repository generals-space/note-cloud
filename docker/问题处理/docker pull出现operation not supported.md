# docker pull出现operation not supported

```
docker pull daocloud.io/centos:6
6: Pulling from centos
32c4f4fef1c6: Extracting [==================================================>] 68.74 MB/68.74 MB
failed to register layer: ApplyLayer exit status 1 stdout:  stderr: symlink gawk /bin/awk: operation not supported
```

情境描述: CentOS7的虚拟机, 编译安装的docker, 版本为`1.14.rc2`, containerd与dockerd服务成功启动, `docker search`也可以正常使用, 但pull操作时出现上述错误.

原因分析: 当时dockerd启动的`graph`参数设置为了从windows宿主机共享的目录, 可能是由于docker需要linux的联合文件系统的特性支持而共享目录本质还是NTFS才会报错.

解决办法: 尝试将`graph`参数指向了linux本地任一目录, 重启dockerd, 再次pull时正常.
