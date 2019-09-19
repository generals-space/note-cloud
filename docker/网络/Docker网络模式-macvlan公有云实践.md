# Docker网络模式-macvlan公有云实践

参考文章

1. [腾讯云 - 开启或关闭广播功能](https://cloud.tencent.com/document/product/215/20116)

2. [弹性网卡在容器网络的应用](http://blog.best-practice.cloud/2019/06/30/%E5%BC%B9%E6%80%A7%E7%BD%91%E5%8D%A1%E5%9C%A8%E5%AE%B9%E5%99%A8%E7%BD%91%E7%BB%9C%E7%9A%84%E5%BA%94%E7%94%A8.html)

最近尝试在阿里云和腾讯云实验了 docker 的 macvlan 网络, 没成功, 容器间无法实现跨主机的网络通信. 

感觉像是公有云提供商将广播包屏蔽掉了, 混杂模式下的网卡也没有办法接收到本来应该发给容器的包.

但是参考文章1中给出了腾讯云中开户子网广播的方法, 试了下, 没用. 难道 macvlan 的网络包不属于广播包?
