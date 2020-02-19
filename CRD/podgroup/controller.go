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

	testgroupv1 "podgroup/pkg/apis/testgroup/v1"
	clientset "podgroup/pkg/client/clientset/versioned"
	crdScheme "podgroup/pkg/client/clientset/versioned/scheme"
	informers "podgroup/pkg/client/informers/externalversions/testgroup/v1"
	listers "podgroup/pkg/client/listers/testgroup/v1"
)

const controllerAgentName = "podgroup-controller"

const (
	// SuccessSynced ...
	SuccessSynced = "Synced"
	// MessageResourceSynced ...
	MessageResourceSynced = "PodGroup synced successfully"
)

// Controller is the controller implementation for PodGroup resources
type Controller struct {
	// kubeclientset is a standard kubernetes clientset
	kubeclientset kubernetes.Interface
	// podGroupClientset is a clientset for our own API group
	podGroupClientset clientset.Interface

	podGroupLister listers.PodGroupLister
	podGroupSynced cache.InformerSynced

	workqueue workqueue.RateLimitingInterface
	recorder record.EventRecorder
}

// NewController returns a new pod group controller
func NewController(
	kubeclientset kubernetes.Interface,
	podGroupClientset clientset.Interface,
	podGroupInformer informers.PodGroupInformer) *Controller {

	utilruntime.Must(crdScheme.AddToScheme(scheme.Scheme))
	glog.V(4).Info("Creating event broadcaster")
	eventBroadcaster := record.NewBroadcaster()
	eventBroadcaster.StartLogging(glog.Infof)
	eventBroadcaster.StartRecordingToSink(&typedcorev1.EventSinkImpl{Interface: kubeclientset.CoreV1().Events("")})
	recorder := eventBroadcaster.NewRecorder(scheme.Scheme, corev1.EventSource{Component: controllerAgentName})

	controller := &Controller{
		kubeclientset:    kubeclientset,
		podGroupClientset: podGroupClientset,
		podGroupLister:   podGroupInformer.Lister(),
		podGroupSynced:   podGroupInformer.Informer().HasSynced,
		workqueue:        workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "PodGroups"),
		recorder:         recorder,
	}

	glog.Info("Setting up event handlers")
	// Set up an event handler for when PodGroup resources change
	podGroupInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: controller.enqueuePodGroup,
		UpdateFunc: func(old, new interface{}) {
			oldPodGroup := old.(*testgroupv1.PodGroup)
			newPodGroup := new.(*testgroupv1.PodGroup)
			if oldPodGroup.ResourceVersion == newPodGroup.ResourceVersion {
                //版本一致，就表示没有实际更新的操作，立即返回
				return
			}
			controller.enqueuePodGroup(new)
		},
		DeleteFunc: controller.enqueuePodGroupForDelete,
	})

	return controller
}

// Run 在此处开始controller的业务
func (c *Controller) Run(threadiness int, stopCh <-chan struct{}) error {
	defer runtime.HandleCrash()
	defer c.workqueue.ShutDown()

	glog.Info("开始controller业务，开始一次缓存数据同步")
	if ok := cache.WaitForCacheSync(stopCh, c.podGroupSynced); !ok {
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
	podGroup, err := c.podGroupLister.PodGroups(namespace).Get(name)
	if err != nil {
		// 如果PodGroup对象被删除了，就会走到这里，所以应该在这里加入执行
		if errors.IsNotFound(err) {
			glog.Infof("PodGroup对象被删除，请在这里执行实际的删除业务: %s/%s ...", namespace, name)

			return nil
		}

		runtime.HandleError(fmt.Errorf("failed to list podGroup by: %s/%s", namespace, name))

		return err
	}

	glog.Infof("这里是podGroup对象的期望状态: %#v ...", podGroup)
	glog.Infof("实际状态是从业务层面得到的，此处应该去的实际状态，与期望状态做对比，并根据差异做出响应(新增或者删除)")

	c.recorder.Event(podGroup, corev1.EventTypeNormal, SuccessSynced, MessageResourceSynced)
	return nil
}

// 数据先放入缓存，再入队列
func (c *Controller) enqueuePodGroup(obj interface{}) {
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
func (c *Controller) enqueuePodGroupForDelete(obj interface{}) {
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
