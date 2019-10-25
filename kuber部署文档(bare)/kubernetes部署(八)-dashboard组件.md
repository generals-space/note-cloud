---
title: kubernetes部署(八)-dashboard组件
tags: [kubernetes, dashboard]
categories: general
---

<!--

# kubernetes部署(八)-dashboard组件

<!tags!>: <!kubernetes!> <!dashboard!>

<!keys!>: tOizwxz79sjns:Gb

-->


参考文章

1. [dashboard官网地址](https://github.com/kubernetes/dashboard)

2. [Kubernetes集群Dashboard插件安装](http://tonybai.com/2017/01/19/install-dashboard-addon-for-k8s/?utm_source=tuicool&utm_medium=referral)

dashboard同dns组件一样, 也是以pods形式存在的, 相关的yaml配置文件是在dashboard主页中提到的[地址](https://git.io/kube-dashboard)下载的.

唯一修改的地方, 就是解开`--apiserver-host=172.32.100.71:8080`的注释, 还是因为没有使用安全接口的问题. 当然, 实际地址以你自己的为准. 这一点是参考文章2中提到的, 与dns插件部署时遇到的问题相似, 而且给出解决方案的是同一位作者[Tony Bai](http://tonybai.com/), 十分感谢.

```yaml
## kube-dashboard.yaml

apiVersion: v1
kind: ServiceAccount
metadata:
  labels:
    k8s-app: kubernetes-dashboard
  name: kubernetes-dashboard
  namespace: kube-system
---
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: ClusterRoleBinding
metadata:
  name: kubernetes-dashboard
  labels:
    k8s-app: kubernetes-dashboard
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: cluster-admin
subjects:
- kind: ServiceAccount
  name: kubernetes-dashboard
  namespace: kube-system
---
kind: Deployment
apiVersion: extensions/v1beta1
metadata:
  labels:
    k8s-app: kubernetes-dashboard
  name: kubernetes-dashboard
  namespace: kube-system
spec:
  replicas: 1
  revisionHistoryLimit: 10
  selector:
    matchLabels:
      k8s-app: kubernetes-dashboard
  template:
    metadata:
      labels:
        k8s-app: kubernetes-dashboard
    spec:
      containers:
      - name: kubernetes-dashboard
        image: gcr.io/google_containers/kubernetes-dashboard-amd64:v1.6.1
        ports:
        - containerPort: 9090
          protocol: TCP
        args:
          # Uncomment the following line to manually specify Kubernetes API server Host
          # If not specified, Dashboard will attempt to auto discover the API server and connect
          # to it. Uncomment only if the default does not work.
          - --apiserver-host=172.32.100.71:8080
        livenessProbe:
          httpGet:
            path: /
            port: 9090
          initialDelaySeconds: 30
          timeoutSeconds: 30
      serviceAccountName: kubernetes-dashboard
      # Comment the following tolerations if Dashboard must not be deployed on master
      tolerations:
      - key: node-role.kubernetes.io/master
        effect: NoSchedule
---
kind: Service
apiVersion: v1
metadata:
  labels:
    k8s-app: kubernetes-dashboard
  name: kubernetes-dashboard
  namespace: kube-system
spec:
  ports:
  - port: 80
    targetPort: 9090
  selector:
    k8s-app: kubernetes-dashboard
```

部署, 也很简单.

```
$ kubectl -s http://172.32.100.71:8080 create -f ./kube-dashboard.yaml 
serviceaccount "kubernetes-dashboard" created
clusterrolebinding "kubernetes-dashboard" created
deployment "kubernetes-dashboard" created
service "kubernetes-dashboard" created
```

insecure模式下通过proxy转发以访问dashboard, 不然没有办法在未认证的情况下连接.

```
$ kubectl -s http://172.32.100.71:8080 proxy --address='0.0.0.0' --port=8001 --accept-hosts='^*$'
Starting to serve on [::]:8001
```

浏览器访问`http://172.32.100.71:8001/ui`

完成.

之后你就需要完成证书认证以及用户名密码登录了, 不然在生产环境这种"裸奔"是绝对不允许的.

> 注意关闭防火墙, 重启docker有可能导致防火墙规则重新生成, 即时清空, 这很重要.