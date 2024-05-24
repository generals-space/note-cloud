# kubectl expose快速创建service[fast]

参考文章

1. [kubectl for Docker Users](https://kubernetes.io/docs/reference/kubectl/docker-cli-to-kubectl/)

最简示例

```
k expose deploy mydeploy --name mysvc --port 8080
```

这会为目标 deploy 创建一个对应的 service 对象. `expose`会自动探测目标 deploy 下面 pod 的 label, 并写到 service 的`spec.selector`中.

可以被 expose 的有 `pod`, `deploy`, 只要指定目标资源的名字即可, 但是不支持`sts`, `ds`.

```log
$ k expose sts 目标sts名称 --name test-svc --port 9092 --target-port=9092
error: cannot expose a StatefulSet.apps
$ k expose ds 目标ds名称 --name test-svc --port 9092 --target-port=9092
error: cannot expose a DaemonSet.extensions
```

## --dry-run -oyaml 快速生成 service 模板

```yaml
## k expose pod esrally(pod名称) --name mysvc --port 8080 --dry-run -oyaml
apiVersion: v1
kind: Service
metadata:
  creationTimestamp: null
  labels:
    app: esrally
    pod-template-hash: cdcf6c456
  name: mysvc
spec:
  ports:
  - port: 8080
    protocol: TCP
    targetPort: 8080
  selector:
    app: esrally
    pod-template-hash: cdcf6c456
status:
  loadBalancer: {}
```
