参考文章

1. [80道kubenetes高频面试题汇总（带答案）](https://blog.51cto.com/yw666/4559012)

rs与rc并没有本质上的不同, 只是名字不一样而已.

既然rs可以单独使用, 为什么还要使用deploy等包裹? 

rs不支持rolling-update, 也不像deploy支持版本记录, 回滚, 暂停升级等高级特性.

docker 创建的ns 无法通过 ip ns 查看到, 那么ta把ns保存在哪了?

openapi 是什么

schema 与 scheme 关于与区别.

shared informer 是什么.

1. api接口设计成`/v1/resources/`的原理
2. client-go库中分布式资源锁的实现原理
3. workqueue 的实现原理

runc 磁盘分配, init 启动命令与信号转发机制

如何手动模拟指定 pod 断网, 然后恢复? 不考虑calico 的网络策略, 只使用 iptables

apiserver watch 接口如何实现.

AdmissionControl（准入机制）是什么, 与认证, 鉴权有什么区别.

list-cache 如何保证数据不丢失

