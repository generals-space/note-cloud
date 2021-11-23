# klog设置日志级别[log level verbosity]

参考文章

1. [How to increase the klog level verbosity?](https://stackoverflow.com/questions/56517463/how-to-increase-the-klog-level-verbosity)

没有找到`klog.SetLevel()`相似的方法, 只能通过从命令行传入`-v`参数.

需要事先将命令行参数绑定到klog中的参数对象.

```golang
import "flag"

klog.InitFlags(flag.CommandLine)
flag.Parse()
```

之后就可以在命令行传入参数了

```
go run main.go --v=5
```
