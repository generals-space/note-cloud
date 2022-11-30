# Pod-健康检查[livenessProbe readinessProbe]

参考文章

1. [Kubernetes Pod 健康检查机制 LivenessProbe 与 ReadinessProbe](https://blog.csdn.net/qq_32641153/article/details/100614499)
    - Pod 的整个生命周期中的状态: `Pending`, `Running`, `Succeeded`, `Failed`, `Unknown`
    - 关于健康检测的部分其实是参考文章2的译文
2. [Configure Liveness, Readiness and Startup Probes](https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-startup-probes/)

```yaml
          livenessProbe:
            exec:
              command:
              - /bin/calico-node
              - -felix-live
            periodSeconds: 10
            initialDelaySeconds: 10
            failureThreshold: 6
          readinessProbe:
            exec:
              command:
              - /bin/calico-node
              - -felix-ready
              - -bird-ready
            periodSeconds: 10
```

`exec`形式的探针的判断依据是命令的退出码, 为0则表示健康检查成功, 否则则表示异常.

```yaml
          livenessProbe:
            failureThreshold: 3
            httpGet:
              path: /healthz
              port: 10254
              scheme: HTTP
            initialDelaySeconds: 10
            periodSeconds: 10
            successThreshold: 1
            timeoutSeconds: 10
          readinessProbe:
            failureThreshold: 3
            httpGet:
              path: /healthz
              port: 10254
              scheme: HTTP
            periodSeconds: 10
            successThreshold: 1
            timeoutSeconds: 10
```

`http`形式的探针的判断依据是http请求的响应状态码, 任何大于或等于200且小于400的代码表示探测成功, 否则则为失败.
