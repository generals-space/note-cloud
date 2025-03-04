# curl发起evict驱逐指令[kubectl]

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

