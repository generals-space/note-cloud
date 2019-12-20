# client-go

参考文章

1. [Authenticating outside the cluster](https://github.com/kubernetes/client-go/tree/master/examples/out-of-cluster-client-configuration)
    - 在集群外, 引用`kubeconfig`作为认证信息.
2. [Authenticating inside the cluster](https://github.com/kubernetes/client-go/tree/master/examples/in-cluster-client-configuration)
    - 在pod内部, 引用`/var/run/secrets/kubernetes.io/serviceaccount`目录内的文件作为认证信息.

## 1. 集群外

引入`k8s.io/client-go/tools/clientcmd`包, 然后使用`BuildConfigFromFlags`得到配置文件对象.

```go
func clientcmd.BuildConfigFromFlags(masterUrl string, kubeconfigPath string) (*rest.Config, error)
```

`masterUrl`一般为"", `kubeconfigPath`即为`~/.kube/config`.

## 2. pod内

引入`k8s.io/client-go/rest`包, 然后使用`InClusterConfig`得到配置文件对象.

```go
func rest.InClusterConfig() (*rest.Config, error)
```

默认读取pod内挂载的`/var/run/secrets/kubernetes.io/serviceaccount`目录下的文件, 生成配置文件对象.

------

获取配置文件对象config后, 再使用`k8s.io/client-go/kubernetes`包下的`NewForConfig`创建`clientset`对象.

```go
func kubernetes.NewForConfig(c *rest.Config) (*kubernetes.Clientset, error)
```

`clientset`拥有各种API, 不同的API拥有不同资源的操作方法.
