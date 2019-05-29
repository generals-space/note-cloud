# docker-容器优雅终止

参考文章

1. [docker容器如何优雅的终止详解](https://www.jb51.net/article/96617.htm)
    - `docker stop`与`docker kill`的信号机制
    - 用golang编写的信号的接收与处理的示例代码
    - CMD指令 列表与字符串两种形式的区别

2. [Sending and Trapping Signals](http://mywiki.wooledge.org/SignalTrap)
    - bash内置命令trap实现信号的捕获与处理
    - trap, sleep, wait

docker与docker-compose都没有类似钩子一样的机制, 无法让容器在启动后/停止前执行指定脚本以完成某种操作, 只有kubernetes中有.

但如果的确有这样的需求, 可以通过docker的`stop`/`kill`子命令实现. 这需要在docker的`CMD`/`ENTRYPOINT`启动命令所运行的进程中捕获信号, 然后做相应的处理. 具体原理可以见参考文章1.

在实际的场景中, docker的CMD命令是`./start.sh && tail -f /etc/profile`. 项目进程由`start.sh`脚本启动, 而实际的pid为1的进程则是`tail`命令. 

项目停止也需要执行对应的`stop.sh`脚本, 之前一直由docker直接停止进程对付过去, 这一次没法对付了...咳.

为了能够捕获docker向程序发出的信号, 需要自行编写信号处理代码. 在实验中我直接使用shell完成, 参考文章1中也提供了golang版本的代码.

`sig.sh`(参考了参考文章2)

```bash
#!/bin/bash

## tailpid初始化
tailpid=
## 设置捕获信号并完成指定操作
trap 'echo yes > /tmp/result.log && kill "$tailpid"' TERM
## &放入后台执行, 类似于开子线程
tail -f /etc/profile &
## 得到上面的job id, 可以看作是线程id
tailpid=$!
## 类似于高级程序语言中的join函数, 阻塞等待子线程完成
wait $tailpid
tailpid=

```

然后将脚本拷贝到镜像中

```dockerfile
## docker build --no-cache=true -f sig.dockerfile -t sigtest .
FROM generals/centos7

COPY sig.sh /root/
CMD /root/sig.sh

## docker run -d --name sigtest -v d:/coding/tmp/sig:/tmp sigtest

```

按照上述命令构建镜像并启动. 然后使用`docker stop sigtest`停止容器, 完成后会发现在`d:/coding/tmp/sig`目录下多了一个`result.log`文件, 其中包含我们在`sig.sh`脚本中输出的内容, 说明我们对信号的处理是成功的.

> 参考文章1中说使用字符串形式的`CMD`指定其实是执行了`bash -c "字符串命令"`, pid为1的是bash进程而不是我们命令中指定的进程. 但是在实验的时候其实两种格式都可以启动为pid=1的进程.

