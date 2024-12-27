# Java容器启动失败-java.lang.OutOfMemoryError：unable to create new native Thread - cgroup：fork rejected by pids controller in system.slice[cgroupfs]

## 问题描述

某个容器的Java进程总是崩溃, 查看日志有"java.lang.OutOfMemoryError：unable to create new native Thread", 最初以为是因为内存资源不足引起的. 

但是同样的容器在另一台主机上运行的好好的, 各种配置都一样, 所以应该不是这个原因.

## 处理思路

查看宿主机`/var/log/message`文件, 有如下输出

```log
2024-12-27T10:26:06.694018+08:00 localhost kernel: [62863.243018] IPv6: ADDRCONF(NETDEV_CHANGE): vethd2f4ea2: link becomes ready
2024-12-27T10:27:03.071960+08:00 localhost kernel: [62919.620455] cgroup: fork rejected by pids controller in /system.slice/docker-36d10fdbeb90705508202f90a61db800f03860d90a6f4ba0bc172a2d9bc7e22f.scope
2024-12-27T10:30:01.808041+08:00 localhost cron[827064]: pam_unix(crond:session): session opened for user root by (uid=0)
2024-12-27T10:30:01.810286+08:00 localhost systemd[1]: Started Session 73 of user root.
2024-12-27T10:30:01.852224+08:00 localhost CRON[827064]: pam_unix(crond:session): session closed for user root
2024-12-27T10:38:10.541931+08:00 localhost containerd[74282]: time="2024-12-27T10:38:10.541827728+08:00" level=info msg="shim disconnected" id=36d10fdbeb90705508202f90a61db800f03860d90a6f4ba0bc172a2d9bc7e22f
2024-12-27T10:38:10.542172+08:00 localhost containerd[74282]: time="2024-12-27T10:38:10.541890786+08:00" level=warning msg="cleaning up after shim disconnected" id=36d10fdbeb90705508202f90a61db800f03860d90a6f4ba0bc172a2d9bc7e22f namespace=moby
2024-12-27T10:38:10.542273+08:00 localhost containerd[74282]: time="2024-12-27T10:38:10.541903790+08:00" level=info msg="cleaning up dead shim"
2024-12-27T10:38:10.542364+08:00 localhost dockerd[143866]: time="2024-12-27T10:38:10.541875318+08:00" level=info msg="ignoring event" container=36d10fdbeb90705508202f90a61db800f03860d90a6f4ba0bc172a2d9bc7e22f module=libcontainerd namespace=moby topic=/tasks/delete type="*events.TaskDelete"
2024-12-27T10:38:10.549244+08:00 localhost containerd[74282]: time="2024-12-27T10:38:10.549172398+08:00" level=warning msg="cleanup warnings time=\"2024-12-27T10:38:10+08:00\" level=info msg=\"starting signal loop\" namespace=moby pid=833301 runtime=io.containerd.runc.v2\n"
2024-12-27T10:38:10.571957+08:00 localhost kernel: [63587.103634] br-315ccae6654f: port 2(vethd2f4ea2) entered disabled state
2024-12-27T10:38:10.571970+08:00 localhost kernel: [63587.103827] veth0aaabbd: renamed from eth0
2024-12-27T10:38:10.659970+08:00 localhost kernel: [63587.191280] br-315ccae6654f: port 2(vethd2f4ea2) entered disabled state
2024-12-27T10:38:10.660004+08:00 localhost kernel: [63587.193274] device vethd2f4ea2 left promiscuous mode
2024-12-27T10:38:10.660006+08:00 localhost kernel: [63587.193277] br-315ccae6654f: port 2(vethd2f4ea2) entered disabled state
2024-12-27T10:38:10.824761+08:00 localhost dockerd[143866]: time="2024-12-27T10:38:10.824677232+08:00" level=warning msg="Your kernel does not support swap limit capabilities or the cgroup is not mounted. Memory limited without swap."
2024-12-27T10:38:10.851566+08:00 localhost systemd-udevd[833401]: Could not generate persistent MAC address for vetha59b1d2: No such file or directory
2024-12-27T10:38:10.851972+08:00 localhost kernel: [63587.384957] br-315ccae6654f: port 2(vethdbc17cd) entered blocking state
2024-12-27T10:38:10.851998+08:00 localhost kernel: [63587.384960] br-315ccae6654f: port 2(vethdbc17cd) entered disabled state
2024-12-27T10:38:10.852000+08:00 localhost kernel: [63587.385171] device vethdbc17cd entered promiscuous mode
2024-12-27T10:38:10.852001+08:00 localhost kernel: [63587.385318] IPv6: ADDRCONF(NETDEV_UP): vethdbc17cd: link is not ready
2024-12-27T10:38:10.852001+08:00 localhost kernel: [63587.385321] br-315ccae6654f: port 2(vethdbc17cd) entered blocking state
2024-12-27T10:38:10.852002+08:00 localhost kernel: [63587.385322] br-315ccae6654f: port 2(vethdbc17cd) entered forwarding state
2024-12-27T10:38:10.853171+08:00 localhost systemd-udevd[833402]: Could not generate persistent MAC address for vethdbc17cd: No such file or directory
2024-12-27T10:38:10.911003+08:00 localhost containerd[74282]: time="2024-12-27T10:38:10.910867351+08:00" level=info msg="loading plugin \"io.containerd.event.v1.publisher\"..." runtime=io.containerd.runc.v2 type=io.containerd.event.v1
2024-12-27T10:38:10.911234+08:00 localhost containerd[74282]: time="2024-12-27T10:38:10.910923290+08:00" level=info msg="loading plugin \"io.containerd.internal.v1.shutdown\"..." runtime=io.containerd.runc.v2 type=io.containerd.internal.v1
2024-12-27T10:38:10.911435+08:00 localhost containerd[74282]: time="2024-12-27T10:38:10.910937455+08:00" level=info msg="loading plugin \"io.containerd.ttrpc.v1.task\"..." runtime=io.containerd.runc.v2 type=io.containerd.ttrpc.v1
2024-12-27T10:38:10.911617+08:00 localhost containerd[74282]: time="2024-12-27T10:38:10.911142550+08:00" level=info msg="starting signal loop" namespace=moby path=/run/containerd/io.containerd.runtime.v2.task/moby/b564ab7252ae20f977d2c8a3e32fd43cb1a75d2866ba8652f6f30518cb8f4514 pid=833436 runtime=io.containerd.runc.v2
2024-12-27T10:38:10.928325+08:00 localhost systemd[1]: Started libcontainer container b564ab7252ae20f977d2c8a3e32fd43cb1a75d2866ba8652f6f30518cb8f4514.
2024-12-27T10:38:11.219956+08:00 localhost kernel: [63587.750212] eth0: renamed from vetha59b1d2
2024-12-27T10:38:11.243942+08:00 localhost kernel: [63587.774247] IPv6: ADDRCONF(NETDEV_CHANGE): vethdbc17cd: link becomes ready
2024-12-27T10:39:14.095973+08:00 localhost kernel: [63650.627009] cgroup: fork rejected by pids controller in /system.slice/docker-b564ab7252ae20f977d2c8a3e32fd43cb1a75d2866ba8652f6f30518cb8f4514.scope
```

"b564ab7252ae20f977d2c8a3e32fd43cb1a75d2866ba8652f6f30518cb8f4514"就是Java容器的pid, 看来真正的原因就在"cgroup: fork rejected by pids controller in /system.slice/"这一行了.

没来得及搜索, 第一眼就怀疑是 systemd cgroup driver 的问题(被坑过). 查看两台宿主机的 docker info, 果然有问题的主机上 cgroup driver 为 systemd.

修改为 cgroupfs 后, 重启 dockerd 服务, 就正常了, good.
