# cgroup添加tasks失败-No space left on device

参考文章

1. [cgroup cpuset No space on the device](https://blog.51cto.com/bingeao/803306)
2. [cgroups: No space left on device](https://forums.fedoraforum.org/showthread.php?286332-cgroups-No-space-left-on-device)
3. [Cgroup change results in "No space left on device" or "Error Code: 5001"](https://access.redhat.com/solutions/232153)

```
$ cgcreate -g cpu,cpuset:gopls
$ cgclassify -g cpuset,cpu:/gopls 4094
Error changing group of pid 4094: No space left on device
```

我自己遇到的是上面的错误, 在搜索解决方法时, 也有如下的错误.

```
echo 2035 > /sys/fs/cgroup/cpuset/gopls/tasks 
-bash: echo: write error: No space left on device
```

实际上两者做的事情是相同的, `cgclassify`会把命令中的pid列表写入`tasks`文件. 而出现问题的原因就在于`cpuset`子系统, 做这些操作前必须设置`cpus`和`mems`两个字段, 其他子系统则没有这个要求.

所以`cgclassify -g cpu:/gopls 4094`不会出现问题, 如果想要把目标pid添加到gopls控制组, 则要先进行如下设置

```
cgset -r cpuset.cpus=0-3 gopls
cgset -r cpuset.mems=0 gopls
```
