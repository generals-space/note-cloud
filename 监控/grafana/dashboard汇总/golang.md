golang 的面板可以直接使用[Go Metrics](https://grafana.com/grafana/dashboards/13722)

ta直接采集prometheus中所有`go_goroutines`指标, 并且已经配置好面板格式. 

创建一个新的 dashboard, 只要把这个面板的 uid 和名称修改一下, 就可以直接使用.

只是数据源, 过滤方式等可能需要自定义.

