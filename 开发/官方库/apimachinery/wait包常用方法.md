```
import (
	"k8s.io/apimachinery/pkg/util/wait"
)
```


`PollImmediateUntil(interval time.Duration, condition ConditionFunc, stopCh <-chan struct{})`

阻塞方法, 定时检测 condition, 直到其返回 true, 或是 stopCh 关闭. 相当于一个轮询方法

`StartWithChannel(stopCh <-chan struct{}, f func(stopCh <-chan struct{}))`

在wait group中启动目标方法, 直到 stopCh 中止.

