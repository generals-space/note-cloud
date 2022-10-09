参考文章

1. [kubernetes wait包详解](https://www.yuque.com/xiaowantang/yu1g2w/iizdgk)

```
import (
	"k8s.io/apimachinery/pkg/util/wait"
)
```

## PollImmediate

```go
PollImmediate(interval, timeout time.Duration, condition ConditionFunc)
```

阻塞方法, 定时检测 condition, 直到其返回 true, 或者超时退出.

## PollImmediateUntil

```go
PollImmediateUntil(interval time.Duration, condition ConditionFunc, stopCh <-chan struct{})
```

阻塞方法, 定时检测 condition, 直到其返回 true, 或是 stopCh 关闭. 相当于一个轮询方法

## StartWithChannel

```go
StartWithChannel(stopCh <-chan struct{}, f func(stopCh <-chan struct{}))
```

在wait group中启动目标方法, 直到 stopCh 中止.

## Jitter

```go
Jitter(duration time.Duration, maxFactor float64)
```

返回一个随机数, duration类型, 用于主调函数进行sleep()

这个随机数在范围在(duration, duration * maxFactor)

- jitter: 抖动
- factor: 因数(乘数)

...这个函数存在的意义是啥? ta底层就是用 rand 做的, 为什么不直接用 rand() 实现, 还要多一层.

## Until

```go
Until(f func(), period time.Duration, stopCh <-chan struct{})
```

每隔一段时间(f), 就执行一次 f 方法, 直到使用 stopCh 结束该过程.
