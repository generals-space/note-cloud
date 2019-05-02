# Linux-locale语系设置

<!--

<!tags!>: <!locale!> <!中文乱码!> <!docker!>

-->

参考文章

1. [CentOS cannot change locale UTF-8解决方法及设置中文支持](http://blog.csdn.net/wave_1102/article/details/45116783)

2. [UTF-8、`en_US.UTF-8`和`zh_CN.UTF-8`的区别](http://blog.csdn.net/huoyunshen88/article/details/41113633)

中文乱码是一个非常常见的问题, 引起中文乱码的原因有很多: 语系, 字体, 编码...其实我也不太清楚一段数据从内存编码到屏幕显示历经了哪些磨难...细想起来感觉很复杂的样子, 水很深, 不想深究.

参考文章2中讲得倒蛮清晰的, 很容易懂.

本文以docker容器中中文乱码为例, 从`locale`入手解决中文乱码的问题.

最初遇到docker中文输入问题时, 进入到容器中输入中文无法上屏. 这是最诡异的情况, 连乱码都没有直接被终端忽略.

然后是这种, 终端输入中文出现半字的情况.

![](https://gitee.com/generals-space/gitimg/raw/master/ff275c714bcd47b80da26730a4a84907.png)

当然, 写入文件也是乱码的.

由于docker镜像过于精简, 所以locale问题应该被最先考虑到. 下面以`centos:6`镜像容器为例.

`locale`: 用于查看系统当前语言设置.

`locale -a`: 查看所有可用语言.

```
[root@db481110976d ~]# locale 
LANG=
LC_CTYPE="POSIX"
LC_NUMERIC="POSIX"
LC_TIME="POSIX"
LC_COLLATE="POSIX"
LC_MONETARY="POSIX"
LC_MESSAGES="POSIX"
LC_PAPER="POSIX"
LC_NAME="POSIX"
LC_ADDRESS="POSIX"
LC_TELEPHONE="POSIX"
LC_MEASUREMENT="POSIX"
LC_IDENTIFICATION="POSIX"
LC_ALL=

[root@db481110976d ~]# locale -a
C
POSIX
en_US.utf8
```

docker容器中所有语言默认设置为`POSIX`, 所以不会支持中文, 而`locale -a`的结果中, 没有可以支持中文的语言. 无论怎样做都不能满足我们的要求.

根据参考文章1中的提示, 我们需要安装`glibc-common`.

```
$ yum install glibc-common -y
```

然后再使用`locale -a`查看容器支持的语言, 你会发现这次的输出有很多, 这里就不列出了.

------

ok, 接下来是如何是我们的语言设置生效的问题.

上面测试中, `locale`的输出中有两个字段是没有值的, `LANG`与`LC_ALL`, 其中直接设置`LC_ALL`的值将会替代所有`LC_*`的值, 而`LANG`则需要单独赋值.

我们执行

```
$ export LANG=zh_CN.UTF-8       ## 有的是zh_CN.UTF-8，不过我在本地没发现这种编码
$ export LC_ALL=zh_CN.UTF-8
```

其实只执行`LANG`的设置命令就可以了, 不会再出现半字的情况, 输出中文到文件也不会乱码了.

> 关于`zh_CN.UTF-8`与`en_US.UTF-8`, 我想区别在于, 设置为前者时, 系统的菜单, 程序的工具栏语言, 输入法默认语言, 时间日期格式, 货币符号等都以中华人民共和国的习惯为准, 而使用后者则会以美国为准. 但两者都同时使用UTF-8字符集, 兼容其他语言的字符而已. --参考文章2.

系统支持的语言存放在`/usr/share/i18n/locales`目录下(以CentOS6为例).

另外, 开机启动的语言设置问题, 除了写入到`/etc/profile`中, 应该还有其他方法. 非docker容器的系统可以尝试写入`/etc/sysconfig/i18n`文件. 不过我没试过, 应该有用.

```
LANG="zh_CN.UTF-8"
```

------

md, 现在系统输出都成中文的了, 真不习惯.

![](https://gitee.com/generals-space/gitimg/raw/master/836c52160fcced1521aaa31068b50ef3.png)