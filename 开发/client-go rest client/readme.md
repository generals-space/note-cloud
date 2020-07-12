参考文章

1. [Kubernetes三种Client的使用示例](https://blog.csdn.net/weiyuanke/article/details/97938690)

kubernetes 的Client库 client-go 中提供了如下三种类型的 client

`ClientSet`: 可以访问集群中所有的原生资源, 如pods、deployment等, 是最常用的一种(无法访问CRD); 

`DynamicClient`: 可以处理集群中所有的资源, 包括crd(自定义资源), 另外它的返回是一个map[string]interface{}类型；目前主要用在garbage collector和namespace controller中; 

`RestClient`: 前面两种client的基础, 更为底层一些; 

本示例的环境:

kuber: 1.16.2
