参考文章

1. [Pod Security Standards](https://kubernetes.io/docs/concepts/security/pod-security-standards/)
2. [Pod 安全性标准](https://kubernetes.io/zh-cn/docs/concepts/security/pod-security-standards/)
    - 参考文章1的中文版
3. [Apply Pod Security Standards at the Namespace Level](https://kubernetes.io/docs/tutorials/security/ns-level-pss/)
    - 通过在 namespace 对象中添加`pod-security.kubernetes.io/enforce=baseline`等标签, 可以在ns级别实现 PSS(Pod安全标准)
4. [Apply Pod Security Standards at the Cluster Level](https://kubernetes.io/docs/tutorials/security/cluster-level-pss/)
    - 在 apiserver 中, 通过`--admission-control-config-file`指定准入控制器的配置文件, 并通过`plugins`字段启用`PodSecurity`控制器, 可以在集群层面实现 PSS(Pod安全标准)

kube: 1.26.0

简单来说, 就是 PSS 定义了3个安全级别, 不满足设置安全级别的行为将被禁止. 这3个安全级别分别为:

- privileged: 不受限制的策略，提供最大可能范围的权限许可。
- baseline: 限制性最弱的策略，禁止已知的策略提升。
- restricted: 限制性非常强的策略，遵循当前的保护 Pod 的最佳实践。

从受限程度上来说, privileged < baseline < restricted.

比如, baseline 策略禁止 hostNetwork, privileged 特权, hostPath 挂载卷等配置; 而 restricted 策略则更进一步, 除了 baseline 中定义的策略外, 还包括禁止以 root 用户运行容器, 必须 drop 所有 capabilities 等要求.

具体的限制策略可见参考文章1和2, 实施方法可见参考文章3和4.

## 使用方法1-命名空间级别

### enforce=baseline

对目标命名空间添加如下标签即可.

```
kubectl label --overwrite ns enforced pod-security.kubernetes.io/enforce=baseline
```

此时创建一个带有 hostPath 挂载卷的 Pod, 根本创建不起来.

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: nginx
  namespace: enforced
spec:
  containers:
  - image: nginx
    name: nginx
    ports:
    - containerPort: 80
    volumeMounts:
    - mountPath: /var/www/html
      name: data
  volumes:
  - name: data
    hostPath:
      path: /opt/nginx
```

```log
$ k apply -f pod.yaml 
Error from server (Forbidden): error when creating "pod.yaml": pods "nginx" is forbidden: violates PodSecurity "baseline:latest": hostPath volumes (volume "data")
```

### warn=baseline

参考文章3和4还提到了`warn=baseline`和`audit=baseline`等几个标签(也可以设置为`restricted`), 可以对目标资源按照`baseline`进行审查, 并发出警告, 但不实际禁止.

```
kubectl label --overwrite ns enforced pod-security.kubernetes.io/enforce-
kubectl label --overwrite ns enforced pod-security.kubernetes.io/warn=baseline
kubectl label --overwrite ns enforced pod-security.kubernetes.io/audit=baseline
```

再次创建上面的 Pod.

```log
$ k apply -f pod.yaml 
Warning: would violate PodSecurity "baseline:latest": hostPath volumes (volume "data")
pod/nginx created
$ k get pod -n enforced
NAME    READY   STATUS    RESTARTS   AGE
nginx   1/1     Running   0          10s
```

虽然有警告, 但实际还是创建成功了.

------

注意, 对已经拥有违反 baseline 策略的 Pod 的 namespace 添加`enforce=baseline`标签是可以成功的, 且不会影响到原本正常运行的 Pod.

```log
$ kubectl label --overwrite ns enforced pod-security.kubernetes.io/enforce=baseline
Warning: existing pods in namespace "enforced" violate the new PodSecurity enforce level "baseline:latest"
Warning: nginx: hostPath volumes
namespace/enforced labeled
$ k get pod -n enforced
NAME    READY   STATUS    RESTARTS   AGE
nginx   1/1     Running   0          4m29s
```

## 使用方法2-集群级别

集群级别的安全标准可以通过在 apiserver 的`--enable-admission-plugins`参数中, 添加`PodSecurity`, 启用该准入控制器(可以说这是一个开关). 

然后通过在`--admission-control-config-file`指定的各准入控制器插件的配置文件中, 添加`PodSecurity`的配置.

```yaml
apiVersion: apiserver.config.k8s.io/v1
kind: AdmissionConfiguration
## plugins 是一个数组
plugins:
- name: PodSecurity
  configuration:
    apiVersion: pod-security.admission.config.k8s.io/v1
    kind: PodSecurityConfiguration
    defaults:
      enforce: "baseline"
      enforce-version: "latest"
      audit: "restricted"
      audit-version: "latest"
      warn: "restricted"
      warn-version: "latest"
    exemptions:
      usernames: []
      runtimeClasses: []
      ## 对 kube-system 命名空间豁免, 忽略该空间下不符合 enforce 策略的 Pod.
      namespaces: [kube-system]
```

