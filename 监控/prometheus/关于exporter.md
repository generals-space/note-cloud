# 关于exporter

参考文章

1. [Prometheus+Grafana搭建监控系统](https://blog.csdn.net/hfut_wowo/article/details/78536022)
    - 以一个初学者的角度解释 prometheus+grafana 监控体系中各组件的作用.

prometheus可以理解为一个数据库+数据抓取工具，工具从各处抓来统一的数据，放入prometheus这一个时间序列数据库中。

那如何保证各处的数据格式是统一的呢？就是通过exporter。exporter也是用GO写的程序，它开放一个http接口，对外提供格式化的数据。所以在不同的环境下，需要编写不同的exporter。

比如, 收集node物理机的信息, 可以使用[node_exptorter](https://github.com/prometheus/node_exporter), 收集tcp, udp, icmp等数据, 可以使用[blackbox_exporter](https://www.github.com/prometheus/blackbox_exporter)等.

每种exporter都可以视作一个服务, 并且通过一个 http metric 接口对外提供数据信息. 比如, `node_exporter`监听的端口默认为9100(因为是监测所有节点的数据, 所以最好通过`DaemonSet`资源进行部署), 访问任意节点上的9100端口, 可以得到如下页面.

![](https://gitee.com/generals-space/gitimg/raw/master/8f87d80804fd4c7f93180f8390742191.png)
