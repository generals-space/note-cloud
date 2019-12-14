参考文章

1. [Installing calicoctl as a Kubernetes pod](https://docs.projectcalico.org/v3.10/getting-started/calicoctl/install#installing-calicoctl-as-a-kubernetes-pod)
    - calicoctl 安装手册
2. [calicoctl user reference](https://docs.projectcalico.org/v3.10/reference/calicoctl/)
    - calicoctl 使用手册

官方提供了3种calicoctl的安装方式: 

1. 二进制文件(裸机部署)
2. docker镜像
3. 作为kuber中的一个pod.

其实这三种方式都算是同一种, 二进制文件需要读取配置文件, 运行的pod也是需要sa等相关权限的支持, 主要目的就是要从etcd或是通过kuber API获取网络状态的信息.


```yaml
---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: calicoctl
  namespace: kube-system

---
apiVersion: v1
kind: Pod
metadata:
  name: calicoctl
  namespace: kube-system
spec:
  nodeSelector:
    beta.kubernetes.io/os: linux
  hostNetwork: true
  ## 官方部署文件缺少hostPID, 结果找不到calico进程
  hostPID: true
  serviceAccountName: calicoctl
  volumes:
    - name: vol-bird
      hostPath:
        path: /var/run/calico
        type: DirectoryOrCreate
  containers:
    - name: calicoctl
      image: calico/ctl:v3.10.1
      command: ["/bin/sh", "-c", "while true; do sleep 3600; done"]
      env:
        - name: DATASTORE_TYPE
          value: kubernetes
      volumeMounts:
        - name: vol-bird
          ## 挂载bird目录
          mountPath: /var/run/bird
---
kind: ClusterRole
apiVersion: rbac.authorization.k8s.io/v1beta1
metadata:
  name: calicoctl
rules:
  - apiGroups: [""]
    resources:
      - namespaces
      - nodes
    verbs:
      - get
      - list
      - update
  - apiGroups: [""]
    resources:
      - nodes/status
    verbs:
      - update
  - apiGroups: [""]
    resources:
      - pods
      - serviceaccounts
    verbs:
      - get
      - list
  - apiGroups: [""]
    resources:
      - pods/status
    verbs:
      - update
  - apiGroups: ["crd.projectcalico.org"]
    resources:
      - bgppeers
      - bgpconfigurations
      - clusterinformations
      - felixconfigurations
      - globalnetworkpolicies
      - globalnetworksets
      - ippools
      - networkpolicies
      - networksets
      - hostendpoints
      - ipamblocks
      - blockaffinities
      - ipamhandles
    verbs:
      - create
      - get
      - list
      - update
      - delete
  - apiGroups: ["networking.k8s.io"]
    resources:
      - networkpolicies
    verbs:
      - get
      - list

---
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: ClusterRoleBinding
metadata:
  name: calicoctl
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: calicoctl
subjects:
- kind: ServiceAccount
  name: calicoctl
  namespace: kube-system
```
