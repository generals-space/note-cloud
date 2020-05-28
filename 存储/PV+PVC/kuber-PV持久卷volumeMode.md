# kuber-PV持久卷volumeMode

参考文章

1. [Volume Mode](https://kubernetes.io/docs/concepts/storage/persistent-volumes/#volume-mode)
    - `volumeMode`字段的意义(一般都使用`Filesystem`类型, 除非显式指定为`Block`)
2. [官方文档 Binding Block Volumes](https://kubernetes.io/docs/concepts/storage/persistent-volumes/#binding-block-volumes)
    - pv/pvc在`volumeMode`的3个可选值: `unspecified(未显式指定)`, `Block`与`Filesystem`之间不同组合时的不同结果(`Bind/No Bind`)
3. [What is File Level Storage vs. Block Level Storage?](https://stonefly.com/resources/what-is-file-level-storage-vs-block-level-storage)
    - `File System`与`Block`的区别

PV和PVC都有`volumeMode`字段, ta们之间的匹配与否关系到pv和pvc实例对象是否能够成功Bind. 

参考文章3介绍了文件存储与块存储的区别. 文件存储可以直接看作对文件/目录操作提供支持, 而块存储则更底层一些, ta可以定义文件系统格式, 一些分布式系统服务, 或是数据库软件可能依赖于特定的文件系统格式, 使用块存储更可靠.
