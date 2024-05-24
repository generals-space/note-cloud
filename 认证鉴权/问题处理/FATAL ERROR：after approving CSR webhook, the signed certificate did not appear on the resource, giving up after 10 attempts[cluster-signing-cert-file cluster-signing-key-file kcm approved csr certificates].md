参考文章

1. [signed certificate did not appear on the resource](https://github.com/kubernetes-sigs/windows-gmsa/issues/19)
2. [Approved Kubernetes CSR, but certificate not shown in status](https://stackoverflow.com/questions/59795325/approved-kubernetes-csr-but-certificate-not-shown-in-status)
    - 高票答案为实际解决方案
3. [一文带你彻底厘清 Kubernetes 中的证书工作机制](https://cloudnative.to/blog/k8s-certificate/)

## 问题描述

在部署 webhook 工程时, 提交 csr 后, 虽然已经是签发状态, 但`webhook-create-signed-cert.sh`脚本仍然执行失败了.

```bash
## 省略...
# approve and fetch the signed certificate
kubectl certificate approve ${csrName}
# verify certificate has been signed
for _ in $(seq 10); do
    serverCert=$(kubectl get csr ${csrName} -o jsonpath='{.status.certificate}')
    if [[ ${serverCert} != '' ]]; then
        break
    fi
    sleep 1
done
if [[ ${serverCert} == '' ]]; then
    ## 这里报错了
    echo "ERROR: After approving csr ${csrName}, the signed certificate did not appear on the resource. Giving up after 10 attempts." >&2
    exit 1
fi
## 省略...
```

虽然 csr 资源的状态已经是`Approved`的, 但是如下命令得到的仍然为空.

```
kubectl get csr ${csrName} -o jsonpath='{.status.certificate}'
```

直接输出yaml, 发现 status 部分为

```yaml
status:
  conditions:
  - lastUpdateTime: "2020-01-17T20:17:20Z"
    message: This CSR was approved by kubectl certificate approve.
    reason: KubectlApprove
    type: Approved
```

## 解决方案

按照参考文章2中高票答案所说, 是因为 controller-manager 中缺少了2个参数

1. --cluster-signing-cert-file 
2. --cluster-signing-key-file

Kubernetes 提供了一个 certificates.k8s.io API，可以使用配置的 CA 根证书来签发用户证书。该 API 由 kube-controller-manager 实现，其签发证书使用的根证书在下面的命令行中进行配置。我们希望 Kubernetes 采用集群根 CA 来签发用户证书，因此在 kube-controller-manager 的命令行参数中将相关参数配置为了集群根 CA。

正常状态下 controller-manager 的参数应该为

```log
$ kubectl get pod -n kube-system kube-controller-manager-hua-dlzx1-i1108-gyt -oyaml | grep sign
    - --cluster-signing-cert-file=/etc/kubernetes/pki/ca.crt # 用于签发证书的 CA 根证书
    - --cluster-signing-key-file=/etc/kubernetes/pki/ca.key # 用于签发证书的 CA 根证书的私钥
    - --cluster-signing-duration=900000h0m0s
    - --controllers=*,bootstrapsigner,tokencleaner
```

添加完这2个参数后, 再重启进程即可.

> 注意: 所有 controller-manager 实例都要添加.

