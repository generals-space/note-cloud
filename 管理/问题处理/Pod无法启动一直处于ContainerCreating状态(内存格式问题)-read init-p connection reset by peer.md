# Pod无法启动一直处于ContainerCreating状态-read init-p connection reset by peer

参考文章

1. [Kubernetes因限制内存配置引发的错误](https://cloud.tencent.com/developer/article/1411527)

场景描述

基本与参考文章1中所说的情况相同, 一直卡在ContainerCreating的状态.

这是因为我的资源配置中, 将内存信息直接写成了`2`, 而不是`2Gi`, 改了就好了.

