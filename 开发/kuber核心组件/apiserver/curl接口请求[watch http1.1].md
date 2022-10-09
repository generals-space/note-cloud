apimachinery/pkg/apis/meta/internalversion/types.go -> ListOptions{}

- limit: int
- resourceVersion: int
- timeoutSeconds

## 常规获取 pod 列表

```
curl -k --max-time 3600 -H 'Content-Type:application/json' --cacert /etc/kubernetes/pki/ca.crt --cert /etc/kubernetes/pki/apiserver-kubelet-client.crt --key /etc/kubernetes/pki/apiserver-kubelet-client.key 'https://127.0.0.1:6443/api/v1/namespaces/default/pods?timeoutSeconds=3600'
```

## pod 资源的 watch 方式

```json
// curl -k --max-time 3600 -H 'Content-Type:application/json' --cacert /etc/kubernetes/pki/ca.crt --cert /etc/kubernetes/pki/apiserver-kubelet-client.crt --key /etc/kubernetes/pki/apiserver-kubelet-client.key 'https://127.0.0.1:6443/api/v1/namespaces/default/pods?timeoutSeconds=3600&watch=true'

{
    "type":"ADDED",
    "object":{
        "kind":"Pod",
        "apiVersion":"v1",
        "metadata":{
            "name":"test-sts-0","generateName":"test-sts-","namespace":"default",
            "selfLink":"","uid":"","resourceVersion":"","creationTimestamp":"",
            "labels":{},
            "ownerReferences":[]
        },
        "spec":{
            "volumes":[{"name":"default-token-jfxlr","secret":{"secretName":"default-token-jfxlr","defaultMode":420}}],
            "volumeMounts":[{"name":"default-token-jfxlr","readOnly":true,"mountPath":"/var/run/secrets/kubernetes.io/serviceaccount"}],
            "containers":[
                {"name":"centos7","image":"registry.cn-hangzhou.aliyuncs.com/generals-space/centos:7","command":["tail","-f","/etc/os-release"],"imagePullPolicy":"IfNotPresent"}
            ],
            "restartPolicy":"Always","terminationGracePeriodSeconds":30,"dnsPolicy":"ClusterFirst","serviceAccountName":"default","serviceAccount":"default","securityContext":{},"hostname":"test-sts-0","subdomain":"test-sts-svc",
            "schedulerName":"middleware-scheduler",
            "priority":0,"enableServiceLinks":true
        },
        "status":{"phase":"Pending","qosClass":"Burstable"}
    }
}
```

> `object`字段包含该 Pod 对象的所有信息.

