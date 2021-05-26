# Quantity.2.数值与计算[limits requests]

## Quantity.AsApproximateFloat64()

> `Approximate`: 近似值

这个方法会打印出一个数值, 整型/浮点型, 看哪个合适.

不过与`String()`有所不同, `String()`的结果中可能包含`k`, `M`等不同的单位后缀.

而`AsApproximateFloat64()`得到的数值是有固定单位的, `cpu`类型的单位为`C(核数)`, `memory`类型的单位则为`Byte`.

```go
	c1, _ := resource.ParseQuantity("10m")
	log.Printf("%+v\n", c1.AsApproximateFloat64()) // 0.01
	c2, _ := resource.ParseQuantity("1k")
	log.Printf("%+v\n", c2.AsApproximateFloat64()) // 1000
	c3, _ := resource.ParseQuantity("1M")
	log.Printf("%+v\n", c3.AsApproximateFloat64()) // 1e+06
```

```go
	m1, _ := resource.ParseQuantity("10Ki")
	log.Printf("%+v\n", m1.AsApproximateFloat64()) // 10240
	m2, _ := resource.ParseQuantity("1Gi")
	log.Printf("%+v\n", m2.AsApproximateFloat64()) // 1.073741824e+09
```

## Quantity.AsInt64()


