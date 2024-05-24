进入到calico-node容器, 查看其中启动的进程.

```log
$ ps -ef
UID         PID   PPID  C STIME TTY          TIME CMD
root          1      0  0 04:20 ?        00:00:00 /usr/bin/runsvdir -P /etc/service/enabled
root         47      1  0 04:20 ?        00:00:00 runsv felix
root         48      1  0 04:20 ?        00:00:00 runsv bird
root         49      1  0 04:20 ?        00:00:00 runsv bird6
root         50      1  0 04:20 ?        00:00:00 runsv confd
root         51     47 10 04:20 ?        00:05:25 calico-node -felix
root         52     50  0 04:20 ?        00:00:02 calico-node -confd
root        136     48  0 04:20 ?        00:00:02 bird -R -s /var/run/calico/bird.ctl -d -c /etc/calico/confd/config/bird.cfg
root        137     49  0 04:20 ?        00:00:02 bird6 -R -s /var/run/calico/bird6.ctl -d -c /etc/calico/confd/config/bird6.cfg
```

其中`bird`, `brid6`, `calico-node`在`/bin`目录下, `confd`, `felix`则没有可执行的二进制文件存在, 应该是内嵌到`calico-node`工程中去了.

由于bird是用C语言写的, 没有与etcd对接, 所以更新配置需要使用`confd`工具.

容器内部

- `/etc/calico/confd/conf.d`: `confd`本身需要读取的配置文件，每个配置文件告诉`confd`模板文件在什么，最终生成的文件应该放在什么地方，更新时要执行哪些操作等.
- `/etc/calico/confd/config`: 生成的配置文件最终放的目录, 也是`bird`读取配置的地方
- `/etc/calico/confd/templates`: 模板文件，里面包括了很多变量占位符，最终会替换成`etcd`中具体的数据.


宿主机上的路由输出如下

```
$ ip r
default via 192.168.0.1 dev ens34 proto static metric 101 
10.23.36.192/26 via 192.168.0.124 dev ens34 proto bird 
10.23.118.64/26 via 192.168.0.125 dev ens34 proto bird 
172.17.0.0/16 dev docker0 proto kernel scope link src 172.17.0.1 
172.32.0.0/24 dev ens33 proto kernel scope link src 172.32.0.121 metric 100 
192.168.0.0/24 dev ens34 proto kernel scope link src 192.168.0.121 metric 101 
192.168.122.0/24 dev virbr0 proto kernel scope link src 192.168.122.1 
```

其中宿主机IP为`192.168.0.121/24`, 而`10.23.36.192/26`与`10.23.118.64/26`则为`coredns`服务的Pod地址.
