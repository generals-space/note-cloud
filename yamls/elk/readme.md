## kibana使用

当所有工作都完成后, 访问kibana的webUI: http://localhost:5601.

实际场景中, 会有多个节点, 多种服务将日志发送到es, 要选出我们想要的服务, 需要先为其添加索引.

进入主界面后, 点击左侧Management -> Elasticsearch[Index Management], 可以看到如下结果

![](https://gitee.com/generals-space/gitimg/raw/master/e8e363d83423a3071e067f6e7d907553.png)

可以看到是按logstash中配置的`nginx-log-%{+YYYY.MM.dd}`格式来的.

> 注意: logstash 需要有 nginx log 目录的读取权限, kuber 自行创建的共享目录可能权限不正确, 最好手动修改`/var/log/nginx`为`755`.

选择Kibana[Index patterns] -> Create index pattern, 输入`nginx-log-*`忽略日期建立索引.

![](https://gitee.com/generals-space/gitimg/raw/master/ec93e67dc0f55df079ad1e230e49c067.png)

![](https://gitee.com/generals-space/gitimg/raw/master/c864dafe2542b2d8f82ed26935ee9350.png)

然后点击左侧Discover.

![](https://gitee.com/generals-space/gitimg/raw/master/45425262341a7364d64bc52807b27d8a.png)

可以看到我们创建的索引已经出现在左侧, 如果有多个索引, 会有下拉框供用户选择, 表示不同的项目.
