# docker中文乱码

参考文章

1. [linux终端不能输入中文解决方法 ](http://blog.sina.com.cn/s/blog_5c4dd3330100cpmm.html)

2. [在Docker容器bash中输入中文](http://blog.shiqichan.com/Input-Chinese-character-in-docker-bash/)

------

docker容器内的bash无论如何都无法输入中文, 不管是在启动容器时打开bash, 还是以服务形式启动容器后再通过 `nsenter`工具进入容器之后显示的bash. 不管是在什么情况下输入甚至粘贴, 不是出现乱码, 回车无反应甚至根本无法上屏. 而且输出时中文全都是乱码.

尝试在容器内 `/root`家目录新建 `.inputrc`文件, 添加以下内容

```shell
set meta-flag on
set convert-meta off
set input-meta on
set output-meta on
```

重启容器发现可以在bash命令行上输入中文, 但是回车发现与预期结果不同, 而且输出时中文依然是乱码. 尝试设置`locale`, 不管将环境变量LANG设置为 `LANG=en_US.UTF-8`还是 `LANG=zh_CN.UTF-8`都不起作用.

------

真正的解决方法是, 在启动容器时传入 `env`参数

```shell
docker run -i -t ubuntu env LANG=C.UTF-8 /bin/bash
```

或是在Dockerfile文件中写入如下行

```shell
ENV LANG=C.UTF-8
```
