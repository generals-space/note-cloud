# StatefulSet-podManagementPolicy 管理策略[sts]

参考文章

1. []

podManagementPolicy:

1. OrderedReady: 一个一个启动, 在第1个 Pod ready 前不会创建后面的 Pod
2. Parallel: 所有 Pod 一起启动, 有可能因为错误导致一起 CrashBackOff

