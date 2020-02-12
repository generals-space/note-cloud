package main

import (
	"fmt"
	"log"
	"path/filepath"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/util/runtime"

	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

/*
	各种处理函数中, obj即为watch接口响应得到的资源对象.
*/

func onAdd(obj interface{}) {
	node := obj.(*corev1.Node)
	fmt.Println("add a node:", node.Name)
}

// onUpdate // 此处省略 workqueue 的使用
func onUpdate(oldObj interface{}, newObj interface{}) {
	fmt.Println("update occurred ===========")
	fmt.Printf("old object: %+v\n", oldObj)
	fmt.Printf("new object: %+v\n", newObj)
}

func onDelete(obj interface{}) {
	fmt.Println("delete not implemented")
}

func main() {
	home := homedir.HomeDir()
	kubeconfig := filepath.Join(home, ".kube", "config")
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err)
	}
	// 初始化 client
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Panic(err.Error())
	}

	stopper := make(chan struct{})
	defer close(stopper)

	// 初始化 informer
	fmt.Printf("new factory...\n")
	// Shared指的是多个 lister 共享同一个cache, 而且资源的变化会同时通知到cache和listers.
	factory := informers.NewSharedInformerFactory(clientset, 0)
	// nodeInformer 拥有两个方法: Informer, Lister.
	// 其实可以把 Informer 看作是 watch 操作.
	nodeInformer := factory.Core().V1().Nodes()
	informer := nodeInformer.Informer()
	defer runtime.HandleCrash()

	// 启动 informer, list & watch
	go factory.Start(stopper)

	// 从 apiserver 同步某种资源的全部对象, 即 list.
	// 之后就可以使用watch这种资源, 维护这份缓存.
	if !cache.WaitForCacheSync(stopper, informer.HasSynced) {
		errTimeout := fmt.Errorf("Timed out waiting for caches to sync")
		runtime.HandleError(errTimeout)
		return
	}

	// 使用自定义 handler, 处理 watch 响应的各种事件.
	// 具体的维护操作在informer内部执行, 这里挂载的是额外的触发器.
	// 需要注意的是, 在上面的list过程中, 会不断触发onAdd事件, 相当于服务发现了.
	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    onAdd,
		UpdateFunc: onUpdate,
		DeleteFunc: onDelete,
	})

	// 从informer对象创建lister, 不过这里的代码没有特殊的目的,
	// 应该只是展示一下通过informer的接口得到list资源的方法.
	nodeLister := nodeInformer.Lister()
	// 从 lister 中获取所有 items
	nodeList, err := nodeLister.List(labels.Everything())
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("get nodelist: ", nodeList)

	<-stopper
}
