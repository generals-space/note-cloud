# Pod-restart自动异常重启机制[倍数 back-off delay]

参考文章

1. [Container restart policy](https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle/#restart-policy)

> The `restartPolicy` applies to all containers in the Pod. 
>
> `restartPolicy` only refers to restarts of the containers by the kubelet on the same node. After containers in a Pod exit, the kubelet restarts them with an exponential back-off delay (10s, 20s, 40s, …), that is capped at five minutes. 
>
> Once a container has executed for 10 minutes without any problems, the kubelet resets the restart backoff timer for that container.

1. `restartPolicy`在pod中设置, 会应用到所有容器.
2. 容器异常结束后, kubelet会自动重启, 但响应时间是递增的; 第1次异常退出, 立刻重启; 第2次异常退出, 10s后重启; 第3次异常退出, 40s后重启, 按倍数递增; 直到多次尝试重启, 问题仍未解决后, 等到容器异常退出5 mins后再重启;
3. 等到pod正常运行10 mins后, 这个递增的时间会重置;
