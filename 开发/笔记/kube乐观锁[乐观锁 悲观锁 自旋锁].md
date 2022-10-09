# kube乐观锁

参考文章

1. [Kubernetes 请求并发控制与数据一致性（含ResourceVersion、Update、Patch简析）](https://blog.csdn.net/jackxuf/article/details/80084358)
    - 介绍了"悲观锁"和乐观锁的概念及各自的优劣.
    - 乐观锁通常通过增加一个资源版本字段, 来判断请求是否冲突.
2. [踩坑client-go for kubernetes乐观锁](https://blog.csdn.net/xf491698144/article/details/106218379/)
    - ResourceVersion是kubernetes实现乐观锁的方式, ResourceVersion使用的是ETCD中的modifiedIndex值
3. [通俗易懂 悲观锁、乐观锁、可重入锁、自旋锁、偏向锁、轻量/重量级锁、读写锁、各种锁及其Java实现！](https://zhuanlan.zhihu.com/p/71156910)
4. [基于etcd的分布式锁](https://www.cnblogs.com/aganippe/p/16011508.html)
    - 分布式锁的特点, 实现方式及优劣, etcd保证分布式锁的能力
    - golang+etcd实现分布式锁示例代码
5. [etcd分布式乐观锁](https://chunlife.top/2019/04/01/etcd%E5%88%86%E5%B8%83%E5%BC%8F%E4%B9%90%E8%A7%82%E9%94%81/)
    - etcd raft原理, quorum模型介绍

