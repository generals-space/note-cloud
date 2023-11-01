参考文章

1. [Could not write JSON: JsonObject; nested exception is com.fasterxml.jackson.databind.JsonMappingException: JsonObject](https://stackoverflow.com/questions/61169128/could-not-write-json-jsonobject-nested-exception-is-com-fasterxml-jackson-data)

## 问题描述

一个 spring boot 项目, 存在如下接口, 功能是查询并返回一个 kube 资源对象.

```java
@GetMapping("/deployments")
@ResponseBody
public V1Deployment get() {
    // 这里的 deploy 应该是从 kube 集群查询出来的某个对象.
    V1Deployment deploy;
    return deploy;
}
```

但在请求时会出现如下错误

```json
{
    "timestamp": 1694077095818,
    "status": 500,
    "error": "Internal Server Error",
    "exception": "org.springframework.http.converter.HttpMessageNotWritableException",
    "message": "Could not write JSON: Java 8 date/time type `java.time.OffsetDateTime` not supported by default: add Module \"com.fasterxml.jackson.datatype:jackson-datatype-jsr310\" to enable handling; nested exception is com.fasterxml.jackson.databind.exc.InvalidDefinitionException: Java 8 date/time type `java.time.OffsetDateTime` not supported by default: add Module \"com.fasterxml.jackson.datatype:jackson-datatype-jsr310\" to enable handling (through reference chain: io.kubernetes.client.openapi.models.V1Deployment[\"metadata\"]->io.kubernetes.client.openapi.models.V1ObjectMeta[\"creationTimestamp\"])",
    "path": "/deployments"
}
```

报错是在接口`return`之后触发的.

## 

这个问题其实之前遇到过, 就是 jackson 库在对 kube 资源对象进行序列化时, 遇到了不支持的日期格式.

解决方案也很简单, 手动加载一个额外模块就可以.

但那个场景是使用了一个自定义的`JSONUtils`工具库, 因此可以自定义加载模块.

而这里的接口是`return`之后触发的, 也就是说, 问题出在 spring boot 自身.

按照参考文章1中所说, 可以在 application.properties 配置文件中为 spring boot 修改默认的序列化工具.

```
Spring Boot >= 2.3.0.RELEASE: spring.mvc.converters.preferred-json-mapper=gson

Spring Boot < 2.3.0.RELEASE: spring.http.converters.preferred-json-mapper=gson
```
