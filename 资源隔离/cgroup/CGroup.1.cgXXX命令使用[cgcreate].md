# Linux-CGroup

参考文章

1. [【Cgroup】Centos7上面的Cgroup简单实验](https://www.cnblogs.com/easton-wang/p/7656205.html)
    - 介绍了cgroup相关命令的使用方法: cgcreate, cgset, cgexec等.
2. [第 3 章 使用 LIBCGROUP 工具](https://access.redhat.com/documentation/zh-cn/red_hat_enterprise_linux/7/html/resource_management_guide/chap-using_libcgroup_tools)
    - cgXXX工具族的使用方法: cgcreate, cgset, cgget, cgexec, cgclassify等.
3. [cgroup实现cpu绑定和资源使用比例限制](https://blog.51cto.com/emulator/1945735)

## 1. cgroup的增删改查

```
cgcreate -g cpu,cpuset:gopls
```

会分别在`/sys/fs/cgroup`下的`cpu`和`cpuset`两个目录下创建`gopls`目录(因为其实实际上创建了两个group), 该目录下会存在`cpu`和`cpuset`两个subsystem的各选项, 参数都为默认值.

> `cgdelete -g cpu,cpuset:gopls`会删除该控制组, 生成的目录也会被删除.

系统中的控制组可以使用`lscgroup`查看

```log
$ lscgroup | grep gopls
cpuset:/gopls
cpu,cpuacct:/gopls
```

注意格式为`子系统:group名称`, 上面的`cgcreate`实际上创建了两个控制组. 

另外, 控制组貌似拥有目录似的结构, 子级目录可以继承父级目录的配置参数, 而上面我们创建的`gopls`实际上存放在了根目录`/`下, 所以`lscgroup`结果中出现了`/gopls`.

接下来使用`cgset`设置cpu和内存的各项控制指标.

```
cgset -r cpu.cfs_period_us=100000 gopls
cgset -r cpu.cfs_quota_us=100000 gopls
cgset -r cpuset.cpus=0-3 gopls
```

可以使用`cgget`查看指定控制组的配置信息.

```
# cgget -g cpuset:/gopls
$ cgget -g cpuset:gopls
gopls:
...省略
cpuset.effective_cpus: 0-3
cpuset.mems:
cpuset.cpus: 0-3
# cgget -g cpu:/gopls
$ cgget -g cpu:gopls
gopls:
...省略
cpu.rt_period_us: 1000000
cpu.cfs_period_us: 100000
cpu.cfs_quota_us: -1
cpu.shares: 1024
```

设置`gopls`控制组各项限制指标, 会修改`/sys/fs/cgroup/cpusest/gopls/cpus`文件的内容为`0-3`, 也会修改对应的`cfs_period_us`和`cfs_quota_us`文件, 可以比较一下.

## 2. cgroup与进程的绑定

上述命令中对cgroup的增删改查操作, 我们还需要把进程归入cgroup控制组, 这样该进程使用的资源就可以被cgroup中定义的规则限制.

cg命令中有两个与进程绑定相关: `cgexec`和`cgclassify`

`cgexec`可以直接在某cgroup下执行命令, 该命令启动的进程使用的资源会受限于此cgroup规则.

```
cgexec -g cpuset,cpu:/gopls command arguments
```

`cgclassify`可以把一个正在运行中的进程移入某控制组并受其控制.

```
cgclassify -g cpuset,cpu:/gopls pidlist
```

`pidlist`为目标进程pid列表, 以空格分隔.

实际上, `cgclassify`会将命令中的pid列表写入`/sys/fs/cgroup/cpu/gopls/tasks`, 这个文件中的进程会被目标cgroup限制资源.
