参考文章

1. [《Kubernetes设计与实现》 16.2 condition设计约定](https://renhongcai.gitbook.io/kubernetes/di-shi-liu-zhang-api-she-ji-yue-ding/1.2-api_convention_condition)

## 约定三：condition需要控制器第一次处理资源时更新

控制器需要尽快地更新condition状态值（condition.status），即便该状态值为Unknown，这么做的好处是可以让其他组件了解到控制器正在调谐这个资源。

然而，并不是所有的控制器都能遵守这个约定，即控制器并不会报告特定的condition（此时该condition状态值可能为Unknown），可能该condition还无法确定，需要在下一次调谐时决定。此种情况下，外部组件无法读取到特定的condition，可以假设该condition为Unknown。

## 约定五：condition不要定义成状态机

condition需要描述当前资源的确定状态，而不是当前资源状态机中的状态。通俗地讲，condition类型需要使用形容词（如Ready）或过去动词（Succeeded），而不是使用当前运行时（如Deploying）。
