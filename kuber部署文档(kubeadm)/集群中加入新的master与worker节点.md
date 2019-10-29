- ##### 
- 


- #### 如果是不记得，请执行以下命令获取

```
1. 在master节点执行kubeadm token list获取token（注意查看是否过期）

2. 如果没有--discovery-token-ca-cert-hash值，也可以通过以下命令获取
openssl x509 -pubkey -in /etc/kubernetes/pki/ca.crt | openssl rsa -pubin -outform der 2>/dev/null | openssl dgst -sha256 -hex | sed 's/^.* //'
```

- 如果是过期了，需要重新生成

```
1. 执行kubeadm token create --print-join-command，重新生成，重新生成基础的 join 命令（对于添加 master 节点还需要重新生成certificate-key，见下一步）
# 如果是添加 worker 节点，不需要执行这一步，直接使用上面返回的 join 命令加入集群。
2. 使用 kubeadm init phase upload-certs --experimental-upload-certs 重新生成certificate-key
# 添加 master 节点：用上面第1步生成的 join 命令和第2步生成的--certificate-key 值拼接起来执行
```

```
# kubeadm init phase upload-certs --experimental-upload-certs
Flag --experimental-upload-certs has been deprecated, use --upload-certs instead
```

参考

