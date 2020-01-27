# CGroup(二)服务的配置文件

参考文章

1. [centos 6,7 上cgroup资源限制使用举例](https://blog.csdn.net/lanyang123456/article/details/81414198)
    - centos 6和7 cgroup组件(`libcgroup`)的安装命令.
    - `cgconfig.conf`及`cgrules.conf`文件的配置
2. [Centos 6 to Centos 7 cgroups](https://serverfault.com/questions/742752/centos-6-to-centos-7-cgroups)
    - `cgconfig.conf`及`cgrules.conf`文件的配置.
3. [docker cgroup技术之cpu和cpuset](https://www.cnblogs.com/charlieroro/p/10281469.html)
    - 介绍了cgroup中的cpu, cpuset两个子系统各自的作用及两个子系统中比较重要的几个参数的含义和使用方法

## `cgconfig`与`cgred`服务

`cgconfig`服务可以将`/etc/cgconfig.d/*`与`/etc/cgrules.conf`这些配置文件进行解析, 并在`/sys/fs/cgroup/`目录下创建对应的子级目录. 比如

```conf
group gopls {
    cpu {
        cpu.cfs_period_us = 100000;
        cpu.cfs_quota_us = 400000;
    }
}
```

上面的配置会创建`/sys/fs/cgroup/cpu/gopls`目录, 该目录下的`cpu.cfs_period_us`内容为`100000`, `cpu.cfs_quota_us`为`40000`.

该服务重启时会对`cgroup`目录下的控制组目录重写, 不过对于从配置文件中删除的控制组并不会被删除, 需要用户手动操作(尤其是对`group`块重命名时, 旧的`group`目录不会删除).

`cgred`服务应该会对`cgrules.conf`中声明的目标进程与系统中的进程进行匹配, 把符合条件的进程pid填写到`/sys/fs/cgroup/${controllers}/${group}/tasks`文件中(该文件中的进程会被所属的`group`控制).

## 字段含义

`cpu.cfs_period_us`: 调度周期, 及将所有待运行任务(有的时候进程处于sleep状态)全部执行一遍的总时间(调度器会根据各任务的权重把这个时间分给各个任务).
`cpu.cfs_quota_us`: 此进程在调度周期当中可以占用的时间

`cpuset.cpus`: 可以指定所属cgroup使用哪个cpu核心. 可用核心可以直接查看`/sys/fs/cgroup/cpuset/cpuset.cpus`文件, 我的服务器是8核, 所以其内容为`0-7`. 该值可以指定为`0`, `0,1`, `0-3`等.

`tasks`: 受该cgroup限制的进程pid列表, 一行一个. 实际上`cgclassify`命令的作用就是把参数中的pid添加到这个文件中.
