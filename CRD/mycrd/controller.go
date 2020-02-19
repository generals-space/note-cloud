package main

import (
	"fmt"
	"time"

	"github.com/golang/glog"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/util/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/util/workqueue"

	mycrdgroupv1 "mycrd/pkg/apis/mycrdgroup/v1"
	clientset "mycrd/pkg/client/clientset/versioned"
	mycrdscheme "mycrd/pkg/client/clientset/versioned/scheme"
	informers "mycrd/pkg/client/informers/externalversions/mycrdgroup/v1"
	listers "mycrd/pkg/client/listers/mycrdgroup/v1"
)

const controllerAgentName = "mycrd-controller"

const (
	// SuccessSynced ...
	SuccessSynced = "Synced"
	// MessageResourceSynced ...
	MessageResourceSynced = "MyCrd synced successfully"
)

// Controller is the controller implementation for MyCrd resources
type Controller struct {
	// kubeclientset is a standard kubernetes clientset
	kubeclientset kubernetes.Interface
	// mycrdClientset is a clientset for our own API group
	mycrdClientset clientset.Interface

	myCrdLister listers.MyCrdLister
	myCrdSynced cache.InformerSynced

	workqueue workqueue.RateLimitingInterface
	recorder record.EventRecorder
}

// NewController returns a new mycrd controller
func NewController(
	kubeclientset kubernetes.Interface,
	mycrdClientset clientset.Interface,
	myCrdInformer informers.MyCrdInformer) *Controller {

	utilruntime.Must(mycrdscheme.AddToScheme(scheme.Scheme))
	glog.V(4).Info("Creating event broadcaster")
	eventBroadcaster := record.NewBroadcaster()
	eventBroadcaster.StartLogging(glog.Infof)
	eventBroadcaster.StartRecordingToSink(&typedcorev1.EventSinkImpl{Interface: kubeclientset.CoreV1().Events("")})
	recorder := eventBroadcaster.NewRecorder(scheme.Scheme, corev1.EventSource{Component: controllerAgentName})

	controller := &Controller{
		kubeclientset:    kubeclientset,
		mycrdClientset: mycrdClientset,
		myCrdLister:   myCrdInformer.Lister(),
		myCrdSynced:   myCrdInformer.Informer().HasSynced,
		workqueue:        workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "MyCrds"),
		recorder:         recorder,
	}

	glog.Info("Setting up event handlers")
	// Set up an event handler for when MyCrd resources change
	myCrdInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: controller.enqueueMyCrd,
		UpdateFunc: func(old, new interface{}) {
			oldMyCrd := old.(*mycrdgroupv1.MyCrd)
			newMyCrd := new.(*mycrdgroupv1.MyCrd)
			if oldMyCrd.ResourceVersion == newMyCrd.ResourceVersion {
                //版本一致，就表示没有实际更新的操作，立即返回
				return
			}
			controller.enqueueMyCrd(new)
		},
		DeleteFunc: controller.enqueueMyCrdForDelete,
	})

	return controller
}

//在此处开始controller的业务
func (c *Controller) Run(threadiness int, stopCh <-chan struct{}) error {
	defer runtime.HandleCrash()
	defer c.workqueue.ShutDown()

	glog.Info("开始controller业务，开始一次缓存数据同步")
	if ok := cache.WaitForCacheSync(stopCh, c.myCrdSynced); !ok {
		return fmt.Errorf("failed to wait for caches to sync")
	}

	glog.Info("worker启动")
	for i := 0; i < threadiness; i++ {
		go wait.Until(c.runWorker, time.Second, stopCh)
	}

	glog.Info("worker已经启动")
	<-stopCh
	glog.Info("worker已经结束")

	return nil
}

func (c *Controller) runWorker() {
	for c.processNextWorkItem() {
	}
}

// 取数据处理
func (c *Controller) processNextWorkItem() bool {

	obj, shutdown := c.workqueue.Get()

	if shutdown {
		return false
	}

	// We wrap this block in a func so we can defer c.workqueue.Done.
	err := func(obj interface{}) error {
		defer c.workqueue.Done(obj)
		var key string
		var ok bool

		if key, ok = obj.(string); !ok {
			c.workqueue.Forget(obj)
			runtime.HandleError(fmt.Errorf("expected string in workqueue but got %#v", obj))
			return nil
		}
		// 在syncHandler中处理业务
		if err := c.syncHandler(key); err != nil {
			return fmt.Errorf("error syncing '%s': %s", key, err.Error())
		}

		c.workqueue.Forget(obj)
		glog.Infof("Successfully synced '%s'", key)
		return nil
	}(obj)

	if err != nil {
		runtime.HandleError(err)
		return true
	}

	return true
}

// 处理
func (c *Controller) syncHandler(key string) error {
	// Convert the namespace/name string into a distinct namespace and name
	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		runtime.HandleError(fmt.Errorf("invalid resource key: %s", key))
		return nil
	}

	// 从缓存中取对象
	mycrd, err := c.myCrdLister.MyCrds(namespace).Get(name)
	if err != nil {
		// 如果MyCrd对象被删除了，就会走到这里，所以应该在这里加入执行
		if errors.IsNotFound(err) {
			glog.Infof("MyCrd对象被删除，请在这里执行实际的删除业务: %s/%s ...", namespace, name)
			return nil
		}

		runtime.HandleError(fmt.Errorf("failed to list mycrd by: %s/%s", namespace, name))

		return err
	}

	glog.Infof("这里是mycrd对象的期望状态: %#v ...", mycrd)
	glog.Infof("实际状态是从业务层面得到的，此处应该去的实际状态，与期望状态做对比，并根据差异做出响应(新增或者删除)")

	c.recorder.Event(mycrd, corev1.EventTypeNormal, SuccessSynced, MessageResourceSynced)
	return nil
}

// 数据先放入缓存，再入队列
func (c *Controller) enqueueMyCrd(obj interface{}) {
	var key string
	var err error
	// 将对象放入缓存
	if key, err = cache.MetaNamespaceKeyFunc(obj); err != nil {
		runtime.HandleError(err)
		return
	}

	// 将key放入队列
	c.workqueue.AddRateLimited(key)
}

// 删除操作
func (c *Controller) enqueueMyCrdForDelete(obj interface{}) {
	var key string
	var err error
	// 从缓存中删除指定对象
	key, err = cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
	if err != nil {
		runtime.HandleError(err)
		return
	}
	//再将key放入队列
	c.workqueue.AddRateLimited(key)
}
