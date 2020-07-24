# prometheus-告警配置alert-manager

参考文章

1. [prometheus告警模块alertmanager注意事项（QQ邮箱发送告警）](https://www.cnblogs.com/danny-djy/p/11097726.html)
    - qq邮箱的配置方式和注意点: 授权码, 465端口及`smtp_require_tls: false`字段设置
2. [官方文档 - DEFINING RECORDING RULES](https://prometheus.io/docs/prometheus/latest/configuration/recording_rules/)
    - 告警规则的配置手册(类型名称, 判断方法, 严重等级, 附加消息等)
3. [官方文档 CONFIGURATION](https://prometheus.io/docs/alerting/configuration/)
    - 告警的处理方式配置手册(发送媒介(邮件, 微信等), 目标配置, 通知模板等)
4. [AlertManager警报通知 E-mail 微信 模板](https://www.cnblogs.com/elvi/p/11444278.html)
    - 微信告警配置及自定义模板方法
5. [alertmanager配置文件说明](https://www.cnblogs.com/zhaojiedi1992/p/zhaojiedi_liunx_65_prometheus_alertmanager_conf.html)
    - `global.resolve_timeout`的解释

用于prometheus的rules告警规则配置如下

```yaml
groups:
  - name: Instances
    rules:
      - alert: InstanceDown
        expr: up == 0
        for: 5m
        labels:
          severity: page
        annotations:
          description: '{{ $labels.instance }} of job {{ $labels.job }} has been down for more than 5 minutes.'
          summary: 'Instance {{ $labels.instance }} down'
  ## rules end ...

  - name: general-service
    ## 用于检测kuber集群中通用服务的健康状况, 当对某service端口探测失败或是service本身处于down的状态, 并持续5m后触发.
    rules:
      - alert: ProbeFailing
        expr: up {job = 'kubernetes-services'} == 0 or probe_success {job = 'kubernetes-services'} == 0
        for: 5m
        labels: 
          serverity: Warning
        annotations:
          description: '{{ $labels.instance }} of job {{ $labels.job }} has been down for more than 5 minutes.'
          summary: 'Service {{ $labels.instance }} down'
```

prometheus webUI -> Alert页面可以见到如下结果

![](https://gitee.com/generals-space/gitimg/raw/master/ffad149afd9850a76beb5bfb984d9d23.jpg)

alert-manager的配置如下

```yaml
global: 
  ## global 一般用来配置邮件发送方(smtp_xxx等), 微信企业版接口认证, webhook的客户端配置(账号密码)等.

  ## 当prometheus发现故障已经解除后并不会立刻发送恢复通知, 
  ## 而是需要再等待`resolve_time`时间进行确认.
  ## 所以在已经修复故障之后, 最多需要等待`global.resolve_timeout`+`receivers.email_configs.send_resolved`的时间.
  ## 但事实并非如此...实际上这个值好像不会影响恢复通知的发送.
  resolve_timeout: 5m
  ## 一般用来配置邮件发送方(smtp_xxx等), 微信企业版接口认证, webhook的客户端配置(账号密码)等.
  ## 邮箱smtp服务器代理
  smtp_smarthost: smtp.qq.com:465
  ## 发送方的邮件地址, smtp服务器以此验证与授权码是否相符
  smtp_auth_username: sender_addr@qq.com
  ## 邮箱密码, 如果是使用qq, 163等, 这里可能需要填写授权码.
  smtp_auth_password: xxxxxxxxx
  ## 邮件发送方的邮箱地址, 也是 email_configs.from 的默认取值.
  smtp_from: sender_addr@qq.com
  ## qq邮箱一般填false
  smtp_require_tls: false
receivers:
  - name: default-receiver
    email_configs:
      ## email_configs中各字段都可以在global块中找到全局配置作为默认值.
      ## 邮件接收方地址, 如果是企业邮箱, 可以通过分组达到群发的目的.
      - to: receiver_addr@qq.com
        ## 故障恢复后是否通知.
        ## 注意: 故障恢复通知需要在`group_interval`时间前确实有警告发生, 
        ## 且之后`group_interval`时间内没有再发生, 才会发送.
        ## (否则alert-manager重启一次启不是要给所有规则都发一次恢复邮件?)
        ## 但也不是立刻发送通知, 而是需要再继续观察`global.resolve_timeout`进行确认.
        send_resolved: true
route:
  receiver: default-receiver
  ## 告警邮件(或是其他类型的通知)发送的频率, 这里是3m一封.
  group_interval: 3m
  ## 最初(即第一次)等待多久时间发送一组警报的通知
  group_wait: 10s
  ## 发送警报的周期
  repeat_interval: 1m
```

当service的pod无法正常运行导致检测失败时, receiver就会收到邮件, 而且是每隔3m一封.

![](https://gitee.com/generals-space/gitimg/raw/master/fe314e1a09dbaa430aa1e79a6d690585.jpg)

![](https://gitee.com/generals-space/gitimg/raw/master/4b065788e07425f657ed45e5875bc0d6.jpg)

> 注意: 邮件内容中的Labels与上面prometheus Alert页面表格中的`Labels`完全相同.

> 默认故障恢复不会有通知信息, 可以通过设置`receivers.email_configs.send_resolved`为true添加. 另外, 虽然官方文档有提到`resolve_timeout`的作用, 但好像并没有生效.

## FAQ

### 1. 

alert-manager的日志

```
level=error ts=2019-12-13T16:15:32.479Z caller=notify.go:367 component=dispatcher msg="Error on notify" err="require_tls: true (default), but \"smtp.qq.com:465\" does not advertise the STARTTLS extension"
level=error ts=2019-12-13T16:15:32.479Z caller=dispatch.go:264 component=dispatcher msg="Notify for alerts failed" num_alerts=1 err="require_tls: true (default), but \"smtp.qq.com:465\" does not advertise the STARTTLS extension"
```

参考文章

1. [prometheus告警模块alertmanager注意事项（QQ邮箱发送告警）](https://www.cnblogs.com/danny-djy/p/11097726.html)

使用qq邮箱作为代理服务器发送邮件时, `smtp_require_tls`需要设置为false, 否则就会出现上述问题.

### 2. 

alert-manager的日志

```
level=error ts=2019-12-14T04:27:08.010Z caller=notify.go:367 component=dispatcher msg="Error on notify" err="cancelling notify retry for \"email\" due to unrecoverable error: parsing from addresses: mail: no angle-addr"
level=error ts=2019-12-14T04:27:08.010Z caller=dispatch.go:264 component=dispatcher msg="Notify for alerts failed" num_alerts=1 err="cancelling notify retry for \"email\" due to unrecoverable error: parsing from addresses: mail: no angle-addr"
```

参考文章

1. [mail: missing phrase](https://github.com/prometheus/alertmanager/issues/624)

`email_configs.from`字段不合法时将会出现上述错误, 一般是格式问题, 无法解析为邮箱地址.
