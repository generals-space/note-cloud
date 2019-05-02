# kubectl-管理集群, 命名空间, 上下文

<!--
<!links!>: Pyak5endi8nCjjr[
-->

以下命令都会修改`kubectl`的配置文件(`~/.kube/config`)

## 1. 新增集群

`kubectl config set-cluster 集群名称 --server=集群apiserver地址`

```
kubectl config set-cluster kubernetes-cluster --server=https://192.168.1.128:8080
```

## 2. 新建开发环境, 并设置当前上下文为此环境

```
## 如果不指定目标集群的话, 则默认为当前所处集群
$ kubectl config set-context fordev
Context "fordev" created.

## 注意, 此时的fordev context下还什么都没有
$ kubectl config get-contexts
CURRENT   NAME       CLUSTER      AUTHINFO           NAMESPACE
*         aliyun     kubernetes   kubernetes-admin   lora
          fordev
          minikube   minikube     minikube
```

...好像不对啊

```
$ kubectl create namespace lora2
namespace "lora2" created

$ kubectl config set-context fordev --cluster=kubernetes --user=kubernetes-admin --namespace=lora2
Context "fordev" modified.
generals-MacBook-Pro:~ general$ kubectl config get-contexts
CURRENT   NAME       CLUSTER      AUTHINFO           NAMESPACE
*         aliyun     kubernetes   kubernetes-admin   lora
          fordev     kubernetes   kubernetes-admin   lora2
          minikube   minikube     minikube
```

看来需要先创建`namespace`, 再把这个`namespace`赋值给某个`context`才行.

接下来切换上下文

```
$ kubectl config use-context fordev
Switched to context "fordev".

$ kubectl config get-contexts
CURRENT   NAME       CLUSTER      AUTHINFO           NAMESPACE
          aliyun     kubernetes   kubernetes-admin   lora
*         fordev     kubernetes   kubernetes-admin   lora2
          minikube   minikube     minikube
```

## 3. 在当前上下文创建新的命名空间并切换当前上下行到到该命名空间下.

```
$ kubectl create namespace lora3
namespace "lora3" created

$ kubectl config get-contexts
CURRENT   NAME       CLUSTER      AUTHINFO           NAMESPACE
*         aliyun     kubernetes   kubernetes-admin   lora
          minikube   minikube     minikube

$ kubectl create namespace lora3
namespace "lora3" created
## 在aliyun环境下, 设置`lora3`命名空间为当前的上下文
$ kubectl config set-context aliyun --namespace=lora3
Context "aliyun" modified.
```