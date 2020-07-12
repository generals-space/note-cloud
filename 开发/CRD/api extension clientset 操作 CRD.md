# api extension clientset 操作 CRD

参考文章

1. [kubernetes/apiextensions-apiserver](https://github.com/kubernetes/apiextensions-apiserver)

通过操作 kuber 内置的资源对象, 使用的是通过`client-go`构建的客户端`clientset/restclient`, 但是ta们无法操作CRD类型. 

开发者如果希望在代码中动态创建/查找/删除 CRD 类型的资源, 需要借助`apiextensions-apiserver`工具库来实现.

当然, 这只是CRD资源, 开发者自定义的资源不叫CRD, 应该叫CR. 操作CR则需要使用通过`code-generator`生成的clientset库了.

下面是使用`apiextensions-apiserver`的一小段示例代码.

```go
package main

import (
	"fmt"
	"path/filepath"

	"github.com/golang/glog"
	apieClientset "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	apimMetav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	cgKuber "k8s.io/client-go/kubernetes"
	cgClientcmd "k8s.io/client-go/tools/clientcmd"
	cgHomedir "k8s.io/client-go/util/homedir"
)

func main() {
	// 先尝试从 ~/.kube 目录下获取配置, 如果没有, 则尝试寻找 Pod 内置的认证配置
	var kubeconfig string
	home := cgHomedir.HomeDir()
	kubeconfig = filepath.Join(home, ".kube", "config")
	cfg, err := cgClientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		glog.Fatalf("Error building kubeconfig: %s", err.Error())
	}

	// kubeClient 用于集群内资源操作, crdClient 用于操作 crd 资源本身.
	// 具体区别目前还不清楚, 不过示例中大多都是这么做的.
	kubeClient, err := cgKuber.NewForConfig(cfg)
	if err != nil {
		glog.Fatalf("Error building kubernetes clientset: %s", err.Error())
	}
	deploy, err := kubeClient.AppsV1().Deployments("kube-system").Get("coredns", apimMetav1.GetOptions{})
	fmt.Printf("%+v\n", deploy)

	apieClient, err := apieClientset.NewForConfig(cfg)
	if err != nil {
		glog.Fatalf("Error building apiextensions clientset: %s", err.Error())
	}
	crds, err := apieClient.ApiextensionsV1().CustomResourceDefinitions().List(apimMetav1.ListOptions{})
	fmt.Printf("%+v\n", crds)
}

```
