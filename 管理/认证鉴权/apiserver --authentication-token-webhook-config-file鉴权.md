# apiserver --authentication-token-webhook-config-file鉴权

apiserver 启动参数中添加

```
    - --authentication-token-webhook-config-file=/etc/kubernetes/pki/webhook_config
```

webhook_config文件的内容如下

```yaml
# clusters refers to the remote service.
clusters:
  - name: name-of-remote-authn-service
    cluster:
      server: http://192.168.30.104:30080/rest/users/auth/token
# users refers to the API server's webhook configuration.

# kubeconfig files require a context. Provide one for the API server.
current-context: webhook
contexts:
- context:
    cluster: name-of-remote-authn-service
  name: webhook
```
