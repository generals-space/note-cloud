# prometheus-告警配置alert-manager自定义邮件模板

参考文章

1. [AlertManager警报通知 E-mail 微信 模板](https://www.cnblogs.com/elvi/p/11444278.html)
    - 微信告警配置及自定义模板方法
2. [alertmanager邮件模版](https://segmentfault.com/a/1190000008695463)
3. [[stable/prometheus]: impossible to add custom templates to alertmanager via alertmanagerFiles](https://github.com/helm/charts/issues/5608)
    - 使用helm安装prometheus时, 自定义模板的编写及配置方法
4. [官方文档 NOTIFICATION TEMPLATE REFERENCE](https://prometheus.io/docs/alerting/notifications/)
    - 通知模板的内置标签及可用函数.


## 1. 发件人与收件人的显示名称

alert-manager默认邮件内容的配置如下

```yaml
receivers:
- name: default-receiver
  email_configs:
  - send_resolved: false
    headers:
      From: sender_addr@qq.com
      To: receiver_addr@qq.com
      Subject: '{{ template "email.default.subject" . }}'
    html: '{{ template "email.default.html" . }}'
```

收到的邮件会有如下效果

![](https://gitee.com/generals-space/gitimg/raw/master/4b065788e07425f657ed45e5875bc0d6.jpg)

但是`headers.From`和`headers.To`这两个键是可以随意改动的, 这两个只表示邮件中显示的名称, 实际的发件人与收件人分别由`email_configs.from`与`email_configs.to`指定. 如果配置成如下

```yaml
    receivers:
      - name: default-receiver
        email_configs:
          - to: 2253238252@qq.com
            headers:
              From: 普罗米修斯
              To: 愚蠢的人类
```

邮件的内容则会是

![](https://gitee.com/generals-space/gitimg/raw/master/9d269897a50c6079ccbe2647f579e8dd.jpg)

可以看到, 显示的发件人和收件人与实际的不相符, 邮箱会有提示, 且头像也无法显示.

## 2. 默认模板中的主题与内容

alert-manager使用的默认邮件模板`email.default.subject`和`email.default.html`都被编译进了最终的二进制文件中, 只能在源码中找到(见`alertmanager/template/default.tmpl`). 由于邮件内容包含了整个html页面, 所以这里只以邮件主题为例, 介绍一下参考文章4中官方文档提到的一些内置标签的涵义.

```go
    [{{ .Status | toUpper }}{{ if eq .Status "firing" }}:{{ .Alerts.Firing | len }}{{ end }}] 
    {{ .GroupLabels.SortedPairs.Values | join " " }} 
    {{ if gt (len .CommonLabels) (len .GroupLabels) }}
        ({{ with .CommonLabels.Remove .GroupLabels.Names }}{{ .Values | join " " }}{{ end }})
    {{ end }}
```

> 注意: 多余的空格会表现在渲染出的内容中, 所以自定义模板完成后可能需要将空格和换行都删除.

对应的邮件主题则是

```

[FIRING:1] (ProbeFailing nginx-svc-01 nginx-svc-01.general-test.svc:80 kubernetes-services nginx-svc-01 general-test Warning)
```

我们对应一下

`.Status`: firing (`toUpper`将其转换为大写)

`.Alerts`: 这是个数组, 表示报警列表, 其中成员字段按照官方文档还是不难理解的.

最难理解的是`CommonLabels`和`GroupLabels`. 我做了实验, 大致确认了如下结果. 首先看一下prometheus Alert页面

![](https://gitee.com/generals-space/gitimg/raw/master/51c70bf46d70b798259dcecd549298d5.jpg)

其中, `kubernetes_namespace`和`kubernetes_name`是在prometheus配置中`scrape_configs.relabel_configs`标签中声明的, `apps`则是在alert-manager配置中`groups.[].rules.labels`中声明的. 其他则是prometheus自己声明的.

`CommonLabels`就是上面的全部.

至于`GroupLabels`, 则是需要在alert-manager配置中`route.group_by`中声明的, 作为分组依据的标签列表的内容. 如果没有这个字段, 则`GroupLabels`为空.

所以`CommonLabels`包含了`GroupLabels`的内容, 上面主题模板中就将后者的字段从前者列表中移除了.
