# Quantity.1.输入与输出[limits requests]

## Quantity{} 对象

```yaml
resources:
  limits:
    cpu: 1
    memory: 1Gi
  requests:
    cpu: 100m
    memory: 100Mi
```

上述内容中, 两对`cpu`与`memory`属于相同的类型, 都是`resource.Quantity`.

`k8s.io/apimachinery/pkg/api/resource`包提供了一个`ParseQuantity()`方法, 用于读取 value 字符串, 得到`Quantity`对象.

```go
quantity, err := resource.ParseQuantity("1G")
```

## cpu 与 limit 区分

如果`ParseQuantity()`的参数字符串带有`m`, `k`, `M`, `G`后缀, 或是直接为整型/浮点型数值(最多可识别小数点后3位), 则可以被当作`cpu`资源值;

如果`ParseQuantity()`的参数字符串带有`Ki`, `Mi`, `Gi`, 则可以被当作`memory`资源值;

### cpu 的换算

- 1000m = 1C
- 1     = 1C
- 1.5   = 1.5C
- 1k    = 1000C
- 1M    = 1000000C
- 1G    = 1000000000C

> `C`是指cpu核心数的意思


### memory 的换算

...略

## Quantity{} 怎么用?

其实`Quantity`是为了简化开发者的计算, `resources.ParseQuantity()`, 与`Quantity.String()`一进一出, 输入和输出的都是人类易读的格式.

`String()`的输出也是会带单位的, 但这其中经过了一些换算, ta会尽量打印出最合适的字符串.

如下是`cpu`资源值的示例.

```go
    // 单位进位
	c1, _ := resource.ParseQuantity("10000m")
	log.Printf("%+v\n", c1.String()) // 10
	c2, _ := resource.ParseQuantity("1000")
	log.Printf("%+v\n", c2.String()) // 1k

    // 单位退位
    c3, _ := resource.ParseQuantity("0.1")
	log.Printf("%+v\n", c3.String()) // 100m
	c4, _ := resource.ParseQuantity("0.1k")
	log.Printf("%+v\n", c4.String()) // 100
```

下面是`memory`资源值的示例.

```
	m1, _ := resource.ParseQuantity("1024Mi")
	log.Printf("%+v\n", m1.String()) // 1Gi
```

由于`memory`是以`1024`为单位, 所以`0.1Gi`是没法得到`100Mi`的, 比cpu要复杂一点.
