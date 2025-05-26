# curl发起evict驱逐指令[kubectl]

参考文章

1. [API 发起的驱逐](https://kubernetes.io/zh-cn/docs/concepts/scheduling-eviction/api-eviction/)
    - 驱逐时, 子资源 Eviction 被创建, 并且 Pod 被删除, 类似于发送一个 DELETE 请求到 Pod 地址.

kubectl 没有 evict 指令, 无法直接发起驱逐指令.

kubectl cordon node名称其实就是发起了 evict 事件.

## 驱逐单个 Pod.

```json
curl -X POST "https://192.168.203.253:6443/api/v1/namespaces/kube-system/pods/coredns-xxx/eviction" \
  -H 'Content-Type:application/json' \
  --cacert /etc/kubernetes/pki/ca.crt \
  --cert /etc/kubernetes/pki/apiserver-kubelet-client.crt \
  --key /etc/kubernetes/pki/apiserver-kubelet-client.key \
  --data '{
    "apiVersion": "policy/v1",
    "kind": "Eviction",
    "metadata": {
      "name": "coredns-xxx",
      "namespace": "kube-system"
    }
  }'
```
