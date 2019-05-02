# pod通过service访问它自己

参考文章

1. [A Pod cannot reach itself via Service IP](https://kubernetes.io/docs/tasks/debug-application-cluster/debug-service/#a-pod-cannot-reach-itself-via-service-ip)

问题没解决...明明node节点已经开启`hairpin`模式了, 而且默认就是`promiscuous-bridge`模式.