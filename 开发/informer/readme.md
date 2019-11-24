

1. [理解 K8S 的设计精髓之 List-Watch机制和Informer模块](https://www.jianshu.com/p/234d27d5c1c1)
    - 关于解释Informer机制的必要性与优势的前言部分写得很好
    - 值得阅读


Etcd存储集群的数据信息, apiserver作为统一入口, 任何对数据的操作都必须经过 apiserver. 客户端(kubelet/scheduler/controller-manager)通过 list-watch 监听 apiserver 中资源(pod/rs/rc等等)的 create, update 和 delete 事件, 并针对事件类型调用相应的事件处理函数. 

list-watch有两部分组成, 分别是list和 watch. list 非常好理解, 就是调用资源的list API罗列资源, 基于HTTP短链接实现; watch则是调用资源的watch API监听资源变更事件, 基于HTTP 长链接实现(http1.1的chunked).

Informer 只会调用Kubernetes List 和 Watch两种类型的 API. 

1. Informer在初始化时, 先调用Kubernetes List API 获得某种 resource的全部Object, 缓存在内存中; 
2. 然后, 调用 Watch API 去watch这种resource, 去维护这份缓存; 
3. 最后, Informer就不再调用Kubernetes的任何 API. 

用List/Watch去维护缓存、保持一致性是非常典型的做法, 但令人费解的是, Informer 只在初始化时调用一次List API, 之后完全依赖 Watch API去维护缓存, 没有任何`resync`机制. 

笔者在阅读Informer代码时候, 对这种做法十分不解. 按照多数人思路, 通过 resync机制, 重新List一遍 resource下的所有Object, 可以更好的保证 Informer 缓存和 Kubernetes 中数据的一致性. 

咨询过Google 内部 Kubernetes开发人员之后, 得到的回复是:

> 在 Informer 设计之初, 确实存在一个relist无法去执 resync操作,  但后来被取消了. 原因是现有的这种 List/Watch 机制, 完全能够保证永远不会漏掉任何事件, 因此完全没有必要再添加relist方法去resync informer的缓存. 这种做法也说明了Kubernetes完全信任etcd. 
