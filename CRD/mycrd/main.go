package main

import (
	"path/filepath"
	"time"

	"github.com/golang/glog"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"

	clientset "mycrd/pkg/client/clientset/versioned"
	informers "mycrd/pkg/client/informers/externalversions"
	"mycrd/pkg/signals"
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

	myCrdClient, err := clientset.NewForConfig(cfg)
	if err != nil {
		glog.Fatalf("Error building example clientset: %s", err.Error())
	}

	myCrdInformerFactory := informers.NewSharedInformerFactory(myCrdClient, time.Second*30)

	//得到controller
	controller := NewController(
		kubeClient,
		myCrdClient,
		myCrdInformerFactory.Mycrdgroup().V1().MyCrds(),
	)

	//启动informer
	go myCrdInformerFactory.Start(stopCh)

	//controller开始处理消息
	if err = controller.Run(2, stopCh); err != nil {
		glog.Fatalf("Error running controller: %s", err.Error())
	}
}
