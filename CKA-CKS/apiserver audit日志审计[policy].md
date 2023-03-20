# apiserver audit日志审计

参考文章

1. [kube-apiserver](https://kubernetes.io/docs/tasks/debug/debug-cluster/audit/#log-backend)
    - 官方文档
    - apiserver的参数选项`--audit-xxx`
2. [Auditing](https://kubernetes.io/docs/tasks/debug/debug-cluster/audit/)
    - 官方文档
    - 审计日志规则: None, Metadata, Request, RequestResponse

kube: 1.26.1

## 开启审计日志

修改`/etc/kubernetes/manifests/apiserver.yaml`中的`command`参数选项

- `--audit-log-path`: 生成的审计日志路径; 
    - 注意, 该文件是在容器内部生成的, 如果要落到宿主机, 需要额外配置目录挂载.
- `--audit-log-maxage`: 审计日志文件能够保留的天数
- `--audit-log-maxbackup`: 审计日志保留的文件个数(日志文件会轮替, 如`message.1`, `message.2`等)
- `--audit-log-maxsize`: 每个日志文件的大小限制(单位为M)
- `--audit-policy-file`: 审计日志的策略文件, 定义要记录哪些行为.

```yaml
  volumeMounts:
  - mountPath: /etc/kubernetes/audit-policy.yaml
    name: audit
    readOnly: true
  - mountPath: /var/log/kubernetes/audit/
    name: audit-log
    readOnly: false

volumes:
- name: audit
  hostPath:
    ## 审计日志策略文件
    path: /etc/kubernetes/audit-policy.yaml
    type: File
- name: audit-log
  hostPath:
    ## 审计日志目录
    path: /var/log/kubernetes/audit/
    type: DirectoryOrCreate
```

## audit policy

apiserver 本质是一个 web 服务器, 可记录的信息包含请示头的认证信息, 请求体与响应体, 在 audit policy 中定义了如下递进的规则

- None: 不记录该事件的任何信息;
- Metadata: 只记录该事件的用户信息, 时间戳, 请求的资源类型, verb等信息, 但不包含请求体与响应体; 
- Request: 在记录 metadata 信息之外, 还记录请求体, 但不包含响应体;
- RequestResponse: 同时包含 metadata, reqeust, 还包含响应体;

上述规则匹配的事件对象是分别进行定义的, 以如下规则为例

```yaml
apiVersion: audit.k8s.io/v1
kind: Policy
rules:
- level: Metadata
```

ta没有定义匹配规则, 因为会匹配所有事件, 即记录所有请示的 metadata 信息.

```yaml
apiVersion: audit.k8s.io/v1
kind: Policy
## rules 中的规则是按顺序匹配的
rules:
- level: Request
  resources:
  - group: "" # core API group
    resources: ["configmaps"]
  namespaces: ["kube-system"]

# Log configmap and secret changes in all other namespaces at the Metadata level.
- level: Metadata
  resources:
  - group: "" # core API group
    resources: ["secrets", "configmaps"]
```

上述规则的详细描述为

1. 记录`kube-system`空间下, 关于`configmap`类型请求的`Request`信息;
2. 记录所有ns空间下的`secret`和除了`kube-system`空间下的`configmap`请示的`Metadata`信息;

**在处理一个请求时, 会按顺序匹配`rules`中定义的规则.**
