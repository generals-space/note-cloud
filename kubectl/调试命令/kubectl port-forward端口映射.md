# kubectl port-forward端口映射

端口映射的目标可以是某个Pod, 也可以是Deployment或是Service这些. 其格式为

将请求从宿主机的8888端口转发到指定服务的5000端口(一般是ClusterIP类型的服务).

```
kubectl port-forward svc/svc名称 8888:5000
```
