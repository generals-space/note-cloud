# Docker+OVS-跨主机容器互联实验

参考文章

1. [基于openvswitch的不同宿主机docker容器网络互联](http://lpyyn.iteye.com/blog/2308714)

2. [利用OpenVSwitch构建多主机Docker网络](http://www.open-open.com/lib/view/open1427619868316.html)

3. [Docker 使用 OpenvSwitch 网桥](http://blog.csdn.net/yeasy/article/details/42555431)

想法来源于kubernetes集群中应用使用flannel完成跨主机容器通信的工作, 新版本的docker也增加了跨主机容器互联的功能, 但都隐藏了技术细节, 不知其所以然. 网上有很多教程, 通过OVS + docker来实现, 其原理相同, 正好可以了解一下.