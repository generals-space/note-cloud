# WorkQueue底层原理

参考文章

1. [深入浅出 kubernetes 之 WorkQueue 详解](https://xie.infoq.cn/article/63258ead84821bc3e276de1f7)
    - WorkQueue 队列与普通FIFO相比, 拥有的额外特性: 有序, 去重, 并发, 标记, 通知, 延迟, 限速, metric等.
    - WorkQueue 接口的3种实现: Interface(FIFO队列接口), DelayingInterface(延迟队列接口), RateLimitingInterface(限速队列接口)
    - WorkQueue 接口的3种实现的具体原理, 有图示, 非常清晰, 值得一看.

