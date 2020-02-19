package main

import (
	"time"
	"path/filepath"

	"github.com/golang/glog"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	clientset "podgroup/pkg/client/clientset/versioned"
	informers "podgroup/pkg/client/informers/externalversions"
	"podgroup/pkg/signals"
	"k8s.io/client-go/util/homedir"

)

func main() {
	// 处理信号
	stopCh := signals.SetupSignalHandler()

	var kubeconfig string
	home := homedir.HomeDir()
	kubeconfig = filepath.Join(home, ".kube", "config")
	cfg, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		glog.Fatalf("Error building kubeconfig: %s", err.Error())
	}

	kubeClient, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		glog.Fatalf("Error building kubernetes clientset: %s", err.Error())
	}

	crdClient, err := clientset.NewForConfig(cfg)
	if err != nil {
		glog.Fatalf("Error building example clientset: %s", err.Error())
	}

	crdInformerFactory := informers.NewSharedInformerFactory(crdClient, time.Second*30)

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
