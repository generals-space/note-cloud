# kuber乐观锁

参考文章

1. [Kubernetes 请求并发控制与数据一致性（含ResourceVersion、Update、Patch简析）](https://blog.csdn.net/jackxuf/article/details/80084358)
    - 介绍了"悲观锁"和乐观锁的概念及各自的优劣.
    - 乐观锁通常通过增加一个资源版本字段, 来判断请求是否冲突.
