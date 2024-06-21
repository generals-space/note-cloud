# CGroup(二)服务的配置文件

参考文章

1. [centos 6,7 上cgroup资源限制使用举例](https://blog.csdn.net/lanyang123456/article/details/81414198)
    - centos 6和7 cgroup组件(`libcgroup`)的安装命令.
    - `cgconfig.conf`及`cgrules.conf`文件的配置
2. [Centos 6 to Centos 7 cgroups](https://serverfault.com/questions/742752/centos-6-to-centos-7-cgroups)
    - `cgconfig.conf`及`cgrules.conf`文件的配置.
3. [docker cgroup技术之cpu和cpuset](https://www.cnblogs.com/charlieroro/p/10281469.html)
    - 介绍了cgroup中的cpu, cpuset两个子系统各自的作用及两个子系统中比较重要的几个参数的含义和使用方法

## 1. 引言

系统环境: CentOS 7

cgroup的控制命令cgXXX(cgcreate, cgset等)做的事情其实就是修改`/sys/fs/cgroup`目录下的配置文件树, 但是这些命令不能持久化, 我们需要以配置文件的形式让cgroup服务加载ta们.

linux中cgroup相关的服务有两个:

1. `cgconfig`: 负责资源控制组挂载
2. `cgred`: 负责识别进程, 并将进程添加到指定资源控制组

> 不过使用`systemctl`启动后前者并没有守护进程存在, 应该是只读取下面的配置文件然后立即退出了, 后者则有`cgrulesengd`进程.

配置文件也有两个: 

1. `/etc/cgconfig.conf`及`/etc/cgconfig.d/*`
2. `/etc/cgrules.conf`

使用`man cgconfig.conf`和`man cgrules.conf`可以分别查看这两个文件的内容格式以及配置示例, 阅读一下能有一个大致的了解. 本示例中只包含了其中的一小部分.

## 2. 示例

使用vscode进程golang工程的远程开发时, 代码补全与智能提示都使用了`gopls`这个golang官方工具(其他编辑器对golang的语法提示插件几乎也都是基于此工具的). 但是这个东西在启动vscode及做一些连续的改动后占用的资源暴涨, 耗尽了开发机器的CPU(8核全占), 内存到是没用多少. 但是对该机器的其他操作全被阻塞住, 有时ssh连接还会被强制断掉. 于是考虑手动对其进行一些限制.

新建`/etc/cgconfig.d/gopls.conf`文件.

```conf
group gopls {
    cpu {
        cpu.cfs_period_us = 100000;
        cpu.cfs_quota_us = 400000;
    }
}
```

然后在`/etc/cgrules.conf`中添加如下行

```
root:gopls cpu gopls
```

然后重启`cgconfig`与`cgred`两个服务, 对名为`gopls`这个进程的资源限制就会生效.

## 3. 解释

### 关于`cgconfig.d/gopls.conf`

`cpu.cfs_period_us`: 调度周期, 及将所有待运行任务(有的时候进程处于sleep状态)全部执行一遍的总时间(调度器会根据各任务的权重把这个时间分给各个任务).
`cpu.cfs_quota_us`: 此进程在调度周期当中可以占用的时间

> centos 7使用的就是完全公平调度(CFS). 感觉cgroup是把`cfs_period_us`和`cfs_quota_us`当成权重值来用了...???

按照上面的配置, 就是允许该cgroup中的进程可以完全使用4个核心的CPU, 在使用`top/ps`等命令查看时, `gopls`的CPU使用最大不会超过`400%`.

可以说, `cpu.cfs_quota_us/cpu.cfs_period_us`的比例决定了目标进程能使用的cpu占用率的上限值(当然绝对不会超过核心数啦).

具体解释可以查看参考文章3.

### 关于`cgrules.conf`

该文件有固定格式

```
<user>                   <controllers>       <destination>
<user>:<process name>    <controllers>       <destination>
```

上述配置文件使用的是第2种, 可以理解为, 为目标进程指定cgroup控制组.

其中`process name`就是`top`命令中的`COMMAND`列的值, 没有参数(我觉得脚本语言在执行时的进程名应该是`bash`, `node`, `python`这种, 与脚本文件名称无关, 因为文件名才是参数).

`controllers`是各`subsystem`的集合, 因为一个`group`块可以声明多个`subsystem`配置, 所以这里也可以写多个, 以逗号分隔.

`destination`则是`group`名称.

## 4. 扩展: 多subsystem的cgroup配置

参考文章1和2中的`group`块都只有一种subsystem, 如果我想在限制进程`gopls`的cpu使用率的同时还想设置其对cpu核心的亲和性(只允许该进程使用某几个cpu)呢?

`cgconfig.d/gopls.conf`配置

```conf
group gopls {
    cpu {
        cpu.cfs_period_us = 100000;
        cpu.cfs_quota_us = 400000;
    }
    cpuset {
        cpuset.cpus = 0-3;
        cpuset.mems = 0;
    }
}
```

`cgrules.conf`配置

```
root:gopls cpu,cpuset gopls
```

重启`cgconfig`和`cgred`服务生效.

`cpuset`就是核心集配置, 你可以使用如下命令查看自己机器上的cpu核心数量

```log
$ cat /sys/fs/cgroup/cpuset/cpuset.cpus
0-7
```

可以使用`1`, `1,3`, `0-3`这种模式指定该cgroup下的进程使用哪些核心.

------

**比较**

在设置`cpuset`块之前, 可以看到每个核心都使用了50%左右(不过也不是每次都这么平均啦).

```log
%Cpu0  :  0.8 us, 49.4 sy,  0.0 ni, 49.4 id,  0.0 wa,  0.0 hi,  0.4 si,  0.0 st
%Cpu1  :  0.4 us, 51.1 sy,  0.0 ni, 48.6 id,  0.0 wa,  0.0 hi,  0.0 si,  0.0 st
%Cpu2  :  0.0 us, 50.7 sy,  0.0 ni, 49.3 id,  0.0 wa,  0.0 hi,  0.0 si,  0.0 st
%Cpu3  :  0.0 us, 51.3 sy,  0.0 ni, 48.7 id,  0.0 wa,  0.0 hi,  0.0 si,  0.0 st
%Cpu4  :  0.4 us, 48.5 sy,  0.0 ni, 51.1 id,  0.0 wa,  0.0 hi,  0.0 si,  0.0 st
%Cpu5  :  0.0 us, 52.0 sy,  0.0 ni, 48.0 id,  0.0 wa,  0.0 hi,  0.0 si,  0.0 st
%Cpu6  :  0.7 us, 54.2 sy,  0.0 ni, 45.1 id,  0.0 wa,  0.0 hi,  0.0 si,  0.0 st
%Cpu7  :  0.4 us, 49.0 sy,  0.0 ni, 50.6 id,  0.0 wa,  0.0 hi,  0.0 si,  0.0 st
```

设置后可以明显看到运算压力多被分配到了`0-3`号核心上.

```log
%Cpu0  : 23.9 us, 62.1 sy,  0.0 ni,  0.7 id,  8.0 wa,  0.0 hi,  5.3 si,  0.0 st
%Cpu1  : 31.5 us, 59.2 sy,  0.0 ni,  0.3 id,  5.1 wa,  0.0 hi,  3.8 si,  0.0 st
%Cpu2  : 28.8 us, 54.3 sy,  0.0 ni,  0.0 id, 11.9 wa,  0.0 hi,  5.0 si,  0.0 st
%Cpu3  : 24.7 us, 66.7 sy,  0.0 ni,  0.0 id,  5.2 wa,  0.0 hi,  3.4 si,  0.0 st
%Cpu4  :  0.0 us,  7.3 sy,  0.0 ni, 52.7 id,  0.0 wa,  0.0 hi,  0.0 si,  0.0 st
%Cpu5  :  0.0 us,  1.7 sy,  0.0 ni, 97.3 id,  0.0 wa,  0.0 hi,  1.0 si,  0.0 st
%Cpu6  :  0.0 us,  0.6 sy,  0.0 ni, 95.5 id,  0.0 wa,  0.0 hi,  3.9 si,  0.0 st
%Cpu7  :  0.0 us,  2.2 sy,  0.0 ni, 46.3 id,  0.0 wa,  0.0 hi,  1.5 si,  0.0 st
```
