# client-go客户端创建

参考文章

1. [Authenticating outside the cluster](https://github.com/kubernetes/client-go/tree/master/examples/out-of-cluster-client-configuration)
    - 在集群外, 引用`kubeconfig`作为认证信息.
2. [Authenticating inside the cluster](https://github.com/kubernetes/client-go/tree/master/examples/in-cluster-client-configuration)
    - 在pod内部, 引用`/var/run/secrets/kubernetes.io/serviceaccount`目录内的文件作为认证信息.
3. [k8s源码解析 - 如何使用yaml创建k8s的资源](https://blog.csdn.net/u014618114/article/details/105168800/)
	- 集群内外访问 apiserver 的客户端的不同创建方式
	- `client-go`中的 clientset, 是各种资源的客户端集合, 实际上操作每一种资源, 都需要相应的客户端, clientSet结构体中每个client成员都是`rest.Interface`类型.

## 1. 获取`rest.Config`配置对象

### 1.1 集群外

引入`k8s.io/client-go/tools/clientcmd`包, 然后使用`BuildConfigFromFlags`得到配置文件对象.

```go
func clientcmd.BuildConfigFromFlags(masterUrl string, kubeconfigPath string) (*rest.Config, error)
```

`masterUrl`一般为"", `kubeconfigPath`即为`~/.kube/config`.

### 1.2 pod内

引入`k8s.io/client-go/rest`包, 然后使用`InClusterConfig`得到配置文件对象.

```go
func rest.InClusterConfig() (*rest.Config, error)
```

默认读取pod内挂载的`/var/run/secrets/kubernetes.io/serviceaccount`目录下的文件, 生成配置文件对象.

## 2. 构建client客户端

获取配置对象后, 有两种方式构建客户端

### 2.1 

使用`k8s.io/client-go/kubernetes`包下的`NewForConfig`创建`clientset`对象.

```go
func kubernetes.NewForConfig(c *rest.Config) (*kubernetes.Clientset, error)
```

`clientset`拥有各种API, 不同的API拥有不同资源的操作方法. 如下

```go
	kubeClient, err := cgKuber.NewForConfig(cfg)
	if err != nil {
		glog.Fatalf("Error building kubernetes clientset: %s", err.Error())
	}
	deploy, err := kubeClient.AppsV1().Deployments("kube-system").Get("coredns", apimMetav1.GetOptions{})
	fmt.Printf("%+v\n", deploy)

```

### 2.2

使用`k8s.io/client-go/rest`包下的`RESTClientFor`创建`RESTClient`对象.

```

```

