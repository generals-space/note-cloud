# kuber-亲和性配置.2.[affinity]

最近遇到一个反亲和性相关的问题.

```yaml
  metadata:
    labels:
      affinitytype: es-general
  spec:
    affinity:
      podAntiAffinity:
        requiredDuringSchedulingIgnoredDuringExecution:
        - labelSelector:
            matchExpressions:
            - key: affinitytype
              operator: In
              values:
              - es-general
          topologyKey: kubernetes.io/hostname
```

上述配置的意思是, 同类型的`Pod`不可调度到同一物理节点(强制).

由于某次配置失误, `metadata.labels.affinitytype`字段的值中少了一个中横线, 成了`esgeneral`, 导致两个节点部署到了同一主机.

由于是在生产环境, 公司规定`Pod`不能随意重启, 在修改之前我评估了一个修改方案. 

最初以为`IgnoredDuringExecution`规则是指如果将`affinitytype`修改回`es-general`, 现在的`Pod`不会发生变化, 然后只要把调度到同一主机的多余的`Pod`删除, 就可以在其他主机上重建, 而且如果修改`replicas`副本值后, 也不会再调度到已有`Pod`的节点上了.

但是我想得太简单了. 修改`metadata.labels.affinitytype`后会导致所有`Pod`重启, 毕竟`controller manager`需要保证所有组件的内容一致, `Pod`的反亲和性必须要和所属的`Deployment`, `StatefulSet`等保持一致才行...当然在此期间位于同一主机的`Pod`会被重新调度.

没办法, 只能先修改一部分业务, 在修改`metadata.labels.affinitytype`同时修改更新策略为`OnDelete`, 手动重启各Pod以保证业务正确性. 剩余比较重要的, 则需要与业务部门协调更新时间了.

--------

直到刚才我还以为`IgnoredDuringExecution`仅仅是指`Pod`标签变化时不会导致`Pod`变动, 而在`Deployment`, `StatefulSet`修改时会导致全体`Pod`变动, 不过看了看上一篇写亲和性的文章...`RequiredDuringSchedulingRequiredDuringExecution`还没支持!?

不过反亲和性的`DuringExecution`阶段应该是鸡肋吧, 根本没有任何情况会用到这种规则.
