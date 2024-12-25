# Docker容器中使用GDB调试出错Operation not permitted

参考文章

1. [docker下gdb调试断点不停](https://blog.csdn.net/so_dota_so/article/details/77509530)

目标程序为go编译成的二进制, 编译命令如下

```
go build -gcflags '-N -l' hello.go
```

使用list查看内容, 以及使用breakpoint打断点都是没问题的, 但是在使用r开始执行程序时就出问题了.

```log
(gdb) b 17
Breakpoint 1 at 0x489431: file /root/hello.go, line 17.
(gdb) r
Starting program: /root/hello
warning: Error disabling address space randomization: Operation not permitted
Cannot create process: Operation not permitted
During startup program exited with code 127.
```

使用lldb也是相同的情况

```log
(lldb) r
error: process launch failed: Child ptrace failed.
```

查阅参考文章1找到原因和解决方法

> linux 内核为了安全起见, 采用了`Seccomp(secure computing)`的沙箱机制来保证系统不被破坏. 它能使一个进程进入到一种"安全"运行模式, 该模式下的进程只能调用4种系统调用(system calls), 即`read()`, `write()`, `exit()`和`sigreturn()`, 否则进程便会被终止. 
> 
> docker只有以`--security-opt seccomp=unconfined`的模式启动容器才能利用GDB调试

如下

```
docker run -it --name golang --security-opt seccomp=unconfined generals/golang_src /bin/bash
```
