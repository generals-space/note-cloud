# ulimit认识

原文链接: [Ulimit](https://zhonghua.io/2019/01/15/linux-ulimit/)

参考文章

1. [如何验证 ulimit 中的资源限制？如何查看当前使用量](https://feichashao.com/ulimit_demo/)
2. [你真知道“Too many open files”?](https://mp.weixin.qq.com/s?__biz=MzIxMjAzMDA1MQ==&mid=2648945736&idx=1&sn=9aa7c240408dd84c4f9d48681f1ec18d&chksm=8f5b5344b82cda52d499cb300514d2b89b0fe6080daeb2bfddcfec427b8b02b4fb9eed4c0fab#rd)
    - 关于`/etc/security/limits.conf`, `/proc/sys/fs/nr_open`
    - `sysctl`修改`fs.file-max`
3. [文件句柄？文件描述符？傻傻分不清楚](文件句柄？文件描述符？傻傻分不清楚)
4. [Linux操作系统中的安全机制—能力（capability)](https://wiki.deepin.io/mediawiki/index.php?title=Linux%E6%93%8D%E4%BD%9C%E7%B3%BB%E7%BB%9F%E4%B8%AD%E7%9A%84%E5%AE%89%E5%85%A8%E6%9C%BA%E5%88%B6---%E8%83%BD%E5%8A%9B%EF%BC%88capability)

------

1. limit的设定值是 per-process 的
2. 在 Linux 中，每个普通进程可以调用`getrlimit()`来查看自己的 limits，也可以调用`setrlimit()`来改变自身的`soft limits`
3. 要改变 hard limit, 则需要进程有`CAP_SYS_RESOURCE`权限
进程 fork() 出来的子进程，会继承父进程的 limits 设定
4. `ulimit`是 shell 的**内置命令**。在执行ulimit命令时，其实是 shell 自身调用`getrlimit()`/`setrlimit()`来获取/改变自身的 limits. 当我们在 shell 中执行应用程序时，相应的进程就会继承当前 shell 的 limits 设定
5. shell 的初始 limits 是谁设定的: 通常是`pam_limits`设定的。顾名思义，pam_limits 是一个 PAM 模块，用户登录后，pam_limits 会给用户的 shell 设定在`limits.conf`定义的值

ulimit, limits.conf 和 pam_limits 的关系，大致是这样的：

1. 用户进行登录，触发 pam_limits;
2. pam_limits 读取 limits.conf，相应地设定用户所获得的 shell 的 limits；
3. 用户在 shell 中，可以通过 ulimit 命令，查看或者修改当前 shell 的 limits;
4. 当用户在 shell 中执行程序时，该程序进程会继承 shell 的 limits 值。于是，limits 在进程中生效了

`/etc/security/limits.conf`格式:

```
#  cat /etc/security/limits.conf
（省略若干....）
# End of file
apps soft nofile 65535
apps hard nofile 65535
apps soft nproc 10240
apps hard nproc 10240
```

第一列表示域（domain）,可以使用用户名（root等），组名（以@开头）,通配置*和%，%可以用于%group参数。

第二列表示类型(type),值可以是soft或者hard

第三列表示项目(item),值可以是core, data, fsize, memlock, nofile, rss, stack, cpu, nproc, as, maxlogins, maxsyslogins, priority, locks, msgqueue, nie, rtprio。其中nofile(Number of Open File）就是文件打开数。

第四列表示值

- **软限制可以在程序的进程中自行改变(突破限制)**
- **硬限制则不行(除非程序进程有root权限)**

## linux capability

Linux内核从2.1版开始提供了”能力”机制,替代传统UNIX模型中”特权/非特权”的简单划分,而是通过具体的能力要求,来判定某一个操作是否可以执行。通过这样的限制,可以确保系统中不存在”完全不被限制的”进程

Linux系统中的能力分为两部分，一部分是进程能力，一部分是文件能力，而Linux内核最终检查的是进程能力中的Effective

每个进程有三个和能力有关的位图：

- permitted(P): 它是effective capabilities和Inheritable capability的超集, 表示进程能够使用的能力, 这些能力是被进程自己临时放弃的, 进程放弃没有必要的能力对于提高安全性大有助益
- effective(E): Linux内核真正检查的能力集
- inheritable(I): 表明该进程可以通过execve继承给新进程的能力

### 可执行文件能力:

TODO

### 进程能力调试:

```log
$ cat /proc/<pid>/status

CapInh:	0000000000000000
CapPrm:	0000003fffffffff
CapEff:	0000003fffffffff
CapBnd:	0000003fffffffff 系统的边界能力
```

### 能力边界集:

- 能力边界集(capability bounding set)是系统中所有进程允许保留的能力。如果在能力边界集中不存在某个能力，那么系统中的所有进程都没有这个能力，即使以超级用户权限执行的进程也一样
- root用户可以向能力边界集中写入新的值来修改系统保留的能力。但是要注意，root用户能够从能力边界集中删除能力，却不能再恢复被删除的能力，只有init进程能够添加能力。通常，一个能力如果从能力边界集中被删除，只有系统重新启动才能恢复
- 删除系统中多余的能力对提高系统的安全性是很有好处的

### 能力可视化:

```
/sbin/capsh --decode=0000003fffffffff
```

获取和设置capabilities:

- 系统调用capget(2)和capset(2)，可被用于获取和设置线程自身的capabilities。此外，也可以使用libcap中提供的接口cap_get_proc(3)和cap_set_proc(3)。当然，Permitted集合默认是不能增加新的capabilities的，除非CAP_SETPCAP在Effective集合中

### Linux能力机制的继承:

```log
P'(permitted) = (P(inheritable) & F(inheritable)) |
                     (F(permitted) & cap_bset)              //新进程的permitted有老进程的和新进程的inheritable和可执行文件的permitted及cap_bset运算得到.
P'(effective) = F(effective) ? P'(permitted) : 0            //新进程的effective依赖可执行文件的effective位，使能：和新进程的permitted一样，负责为空
P'(inheritable) = P(inheritable)    [i.e., unchanged]       //新进程的inheritable直接继承老进程的Inheritable
```

说明:

P   在执行execve函数前，进程的能力
P'  在执行execve函数后，进程的能力
F   可执行文件的能力
cap_bset 系统能力的边界值，在此处默认全为1
