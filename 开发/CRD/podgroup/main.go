package main

import (
	"path/filepath"
	"time"

	"github.com/golang/glog"
	cgKuber "k8s.io/client-go/kubernetes"
	cgClientcmd "k8s.io/client-go/tools/clientcmd"
	cgHomedir "k8s.io/client-go/util/homedir"

	clientset "podgroup/pkg/client/clientset/versioned"
	crdInformerFactory "podgroup/pkg/client/informers/externalversions"
	"podgroup/pkg/signals"
)

func main() {
	// 处理信号
	stopCh := signals.SetupSignalHandler()

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
	crdClient, err := clientset.NewForConfig(cfg)
	if err != nil {
		glog.Fatalf("Error building example clientset: %s", err.Error())
	}

	crdInformerFactory := crdInformerFactory.NewSharedInformerFactory(crdClient, time.Second*30)

	//得到controller
	controller := NewController(
		kubeClient,
		crdClient,
		crdInformerFactory.Testgroup().V1().PodGroups(),
	)

	//启动informer
	go crdInformerFactory.Start(stopCh)

	//controller开始处理消息
	if err = controller.Run(2, stopCh); err != nil {
		glog.Fatalf("Error running controller: %s", err.Error())
	}
}
