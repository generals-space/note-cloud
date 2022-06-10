# Expected a string but was BEGIN_OBJECT at line 1 column 1152 path .status

参考文章

1. [Getting exception from deleteNamespacedStatefulSet() call](https://github.com/kubernetes-client/java/issues/86)
2. [deleteCollectionNamespacedJob](https://github.com/kubernetes-client/java/blob/client-java-parent-5.0.0/kubernetes/docs/BatchV1Api.md)

kuber: 1.13.2

java client: 5.0.0

java client 与 kuber 集群的版本是匹配的.

## 场景描述

在使用`BatchV1API`批量删除`Job`时, 出现了如下异常.

```java
    BatchV1Api batchV1Api = kubeclient.getBatchV1Api(xxx);
    try {
        ApiResponset<V1Status> apiResponse = batchV1Api.deleteNamespacedJobWithHttpInfo(
            name, namespace, null, body, null, null, false, null
        );
        if(HttpUtils.isSuccess(apiResponse.getStatusCode())){
            // 
        }
    } catch(Exception e){
        // 每次都运行到这里...
    }
```

```
com.google.gson.JsonSyntaxException: java.lang.IllegalStateException: Expected a string but was BEGIN_OBJECT at line 1 column 1152 path $.status
```

而实际上, 目标资源对象已经被删除了, 只是在返回时对响应的解析失败了而已. 但这样就会导致只有匹配到的第一个Job被删除了, 剩余的都还在...

不只删除 Job, 所有删除操作其实都会报这个异常, 这是这个版本的 Bug. 

## 解决方案

单个删除时可以直接忽略这个接口的返回结果, 默认ta已经成功就行了. 但是我想要的是批量删除的能力, 所以只能对这个类型的异常单独处理一下.

```java
    BatchV1Api batchV1Api = kubeclient.getBatchV1Api(xxx);

    try {
        ApiResponset<V1Status> apiResponse = batchV1Api.deleteNamespacedJobWithHttpInfo(
            name, namespace, null, body, null, null, false, null
        );
        LOGGER.debug("Deleted namespace {}", namespaceString);
    } catch (JsonSyntaxException e) {
        if (e.getCause() instanceof IllegalStateException) {
            IllegalStateException ise = (IllegalStateException) e.getCause();
            if (ise.getMessage() != null && ise.getMessage().contains("Expected a string but was BEGIN_OBJECT"))
                LOGGER.debug("Catching exception because of issue https://github.com/kubernetes-client/java/issues/86", e);
            else 
                throw e;
        }
        else 
            throw e;
    }
```

## deleteCollectionXXX() 批量删除

另外, 参考文章1中还有提到了`deleteCollectionXXX()`方法, 其实本来也是可以的. 但是无法进行关联删除, 因为这个版本的 java client `deleteCollectionXXX()`方法中不包含`orphanDependents`和`propagationPolicy`参数...
