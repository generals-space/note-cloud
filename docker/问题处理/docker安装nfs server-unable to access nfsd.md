# docker安装nfs server-unable to access nfsd

```
docker run -d --name nfs-server --net host --pid host -v /opt/nfs:/opt/nfs --env nfspath="/opt/nfs 192.168.80.0/24" generals/nfs-server
```

```log
Starting rpcbind: [  OK  ]
FATAL: Could not load /lib/modules/3.10.0-1062.4.1.el7.x86_64/modules.dep: No such file or directory
Starting NFS services:  [  OK  ]
Starting NFS mountd: [  OK  ]
rpc.nfsd: Unable to access /proc/fs/nfsd errno 2 (No such file or directory).
Please try, as root, 'mount -t nfsd nfsd /proc/fs/nfsd' and then restart rpc.nfsd to correct the problem
Starting NFS daemon: [FAILED]
/opt/nfs 192.168.80.0/24(rw,sync,no_all_squash)
```

最开始以为是`/proc`目录的问题, 就添加了`--pid`选项, 共享宿主机的目录, 结果还是不行.

后来添加上`--privileged`选项就可以了.
