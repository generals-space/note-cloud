# patch -f 指定文件

patch子命令没有直接使用`-f`等选项指定读取哪个文件, 但是可以在命令行配合`cat`使用.

```json
{
    "spec": {
        "template": {
            "spec": {
                "containers": [
                {
                    "name": "centos7",
                    "image": "registry.cn-hangzhou.aliyuncs.com/generals-space/centos:7-devops"
                }
                ]
            }
        }
    }
}
```

```
k patch ds test-ds --patch "$(cat patch-file.yaml)"
```

> 虽然`patch`不支持从文件读入更改内容, 但是`kubectl`还有一个子命令`replace`可以实现, 看名字就知道, ta的替换类型是全部替换...

