# kubectl port-forward端口映射

> 需要宿主机安装`socat`工具.

端口映射的目标可以是某个Pod, 也可以是Deployment或是Service这些, 需要通过`type/name`这种格式进行区分. 当然, 实际上转发的目标最终还是Pod, 如果你指向了Service或是Deployment, 该命令会选择一个匹配的Pod进行转发, 如果该Pod退出了, 会自动重新选一个再次映射.

将请求从宿主机的8888端口转发到指定服务的5000端口(一般是ClusterIP类型的服务).

## 示例

监听本地的5000和6000端口, 并将其映射到mypod的5000和6000端口

```
kubectl port-forward pod/mypod 5000 6000
```

监听本地的8888端口, 将其映射到mypod的5000端口. 

```
kubectl port-forward pod/mypod 8888:5000
```

不过需要注意的是, 这种格式的8888只监听在`127.0.0.1`的地址上, 你需要添加`--address 0.0.0.0`选项.

```
kubectl port-forward --address 0.0.0.0 pod/mypod 8888:5000
```

另外, 由于Service的端口和Pod内部的端口可以不一致, 所以如果目标对象为Service时, 则需要将目标端口填写成Service暴露出来的端口. 以如下Service为例

```yaml
apiVersion: v1
kind: Service
metadata:
  name: mysvc
spec:
  ports:
  - name: http
    port: 80
    protocol: TCP
    targetPort: 9090
```

如果目标为Pod, 则命令可写作

```
kubectl port-forward --address 0.0.0.0 pod/mypod 8888:9090
```

而如果目标是Service, 则命令应该为

```
kubectl port-forward --address 0.0.0.0 svc/mysvc 8888:80
```
