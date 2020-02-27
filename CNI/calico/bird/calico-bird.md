# calico-bird

参考文章

1. [官方手册 BIRD 2.0 User's Guide (version 2.0.7)](https://bird.network.cz/?get_doc&v=20&f=bird.html#toc6)
    - 目录
2. [官方手册 6. Protocols](https://bird.network.cz/?get_doc&v=20&f=bird-6.html#ss6.3)
    - `template`块配置
    - 某些字段类型为`switch`, 其实就是布尔类型, 取值为`on/off`.
2. [官方手册 3. Configuration](https://bird.network.cz/?get_doc&v=20&f=bird-3.html)
    - `protocol`块配置
    - `export filter`选项解释

BGP协议并不简单, 我们这里先看看bird的使用方法, 对应到calico中起到了什么作用.

在`calico-node`容器中, 启动了`bird -R -s /var/run/calico/bird.ctl -d -c /etc/calico/confd/config/bird.cfg`进程.

- `-c`: 指定配置文件
- `-d`: debug模式, 同时在前台运行.
- `-s`: `control socket`, 这是一个Unix Socket文件.
- `-R`: 启用平滑重启机制.

`bird`进程监听179端口.

------

bird的配置文件`bird.cfg`有点像`nginx.conf`.

