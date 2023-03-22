# java client - io.kubernetes.client.ApiException - Bad Request查看具体原因[setDebugging(true)]

参考文章

1. [How to find out root cause "io.kubernetes.client.ApiException: Bad Request" exception? ](https://github.com/kubernetes-client/java/issues/695)
        - 亲测有效

- kubernetes: v1.23.4
- kubernetes-client-java: 5.0.0

现在正在用的 java client api 可以适配 1.13.x, 1.17.x 的 kubernetes 集群, 不过在与 1.23.4 版本的集群时, 使用 corev1 查询 pvc 资源时, 出现了问题(可是创建pvc是没问题的🤨).

```java
        try {
            ApiResponse<V1PersistentVolumeClaim> apiResponse = coreV1Api.readNamespacedPersistentVolumeClaimWithHttpInfo(pvcName, namespace, null, true, true);
            if (HttpUtils.isSuccess(apiResponse.getStatusCode())) {
                return apiResponse.getData();
            }
        } catch (Exception e) {
            logger.error("get pvc failed ",  e.getMessage());
        }
```

出问题的行是"readNamespacedPersistentVolumeClaimWithHttpInfo()"函数的所在行, 异常被 catch 语句捕获, 如下

```
io.kubernetes.client.ApiException: Bad Request
        at io.kubernetes.client.ApiClient.handleResponse(ApiClient.java:886)
        at io.kubernetes.client.ApiClient.execute(ApiClient.java:802)
        at io.kubernetes.client.apis.CoreV1Api.readNamespacedPersistentVolumeClaimWithHttpInfo(CoreV1Api.java:24668)
        at com.cmos.k8s.middleware.client.impl.PersistentVolumeClaimApisImpl.getPersistentVolumeClaimByNameAndNamespaces(PersistentVolumeClaimApisImpl.java:71)
        at com.cmos.k8s.middleware.service.es.impl.ElasticSearchServiceImpl.updateElasticsearchLocalStorage(ElasticSearchServiceImpl.java:394)
```

...但是这个报错完全看不出问题嘛, "readNamespacedPersistentVolumeClaimWithHttpInfo()"是 jar 包里的函数, apiResponse 变量也得不到就直接被 catch 了.

后来找到参考文章1, 在构建`ApiClient`对象是, 添加一行"client.setDebugging(true);"

```java
    @Override
    public ApiClient buildApiClient(Cluster cluster) {
        String address = String.format("%s://%s:%d", cluster.getProtocol(), cluster.getHost(), cluster.getPort());
        return new ClientBuilder().setBasePath(address).setVerifyingSsl(false)
                .setAuthentication(new AccessTokenAuthentication(cluster.getMachineToken())).build();
    }

    @Override
    public CoreV1Api getCoreV1Api(Cluster cluster) {
        ApiClient apiClient = this.buildApiClient(cluster);
        apiClient.setDebugging(true); // 新增此行
        return new CoreV1Api(apiClient);
    }
```

然后就可以打印出请求 apiserver 的参数, 请求体, 以及响应体信息了.

```
--> GET https://172.22.248.183:6443/api/v1/namespaces/zjjpt-es/persistentvolumeclaims/general-es-0524-01--claim?exact=true&export=true HTTP/1.1
authorization: Bearer 证书信息
Accept: application/json
User-Agent: Swagger-Codegen/1.0-SNAPSHOT/java
--> END GET
<-- HTTP/1.1 400 Bad Request (47ms)
Audit-Id: 048013f6-1967-4ea7-923d-e81d564aeda1
Cache-Control: no-cache, private
Content-Type: application/json
X-Kubernetes-Pf-Flowschema-Uid: 6a01202f-0ea2-432c-950b-930d45789524
X-Kubernetes-Pf-Prioritylevel-Uid: 5fac74b2-ebc2-4fa8-9867-c2d06a48851d
Date: Wed, 25 May 2022 03:11:30 GMT
Content-Length: 183
OkHttp-Sent-Millis: 1653448290278
OkHttp-Received-Millis: 1653448290279

{"kind":"Status","apiVersion":"v1","metadata":{},"status":"Failure","message":"the export parameter, deprecated since v1.14, is no longer supported","reason":"BadRequest","code":400}

<-- END HTTP (183-byte body)
2022-05-25 11:11:30,280 ERROR http-nio-0.0.0.0-18080-exec-2 (com.cmos.k8s.middleware.client.impl.PersistentVolumeClaimApisImpl:76) - get pvc failed  + Bad Request
io.kubernetes.client.ApiException: Bad Request
        at io.kubernetes.client.ApiClient.handleResponse(ApiClient.java:886)
        at io.kubernetes.client.ApiClient.execute(ApiClient.java:802)
        at io.kubernetes.client.apis.CoreV1Api.readNamespacedPersistentVolumeClaimWithHttpInfo(CoreV1Api.java:24668)
        at com.cmos.k8s.middleware.client.impl.PersistentVolumeClaimApisImpl.getPersistentVolumeClaimByNameAndNamespaces(PersistentVolumeClaimApisImpl.java:71)
        at com.cmos.k8s.middleware.service.es.impl.ElasticSearchServiceImpl.updateElasticsearchLocalStorage(ElasticSearchServiceImpl.java:394)
```

可以看到, 是因为 kubernetes 在 1.23.4 版本已经不再接受 exporter 参数了, 所以在调用"readNamespacedPersistentVolumeClaimWithHttpInfo()"时, 需要将`exporter`参数置为 null.

```java
        try {
            // ApiResponse<V1PersistentVolumeClaim> apiResponse = coreV1Api.readNamespacedPersistentVolumeClaimWithHttpInfo(pvcName, namespace, null, true, true);
            // the export parameter, deprecated since v1.14, is no longer supported
            ApiResponse<V1PersistentVolumeClaim> apiResponse = coreV1Api.readNamespacedPersistentVolumeClaimWithHttpInfo(pvcName, namespace, null, true, null);
            if (HttpUtils.isSuccess(apiResponse.getStatusCode())) {
                return apiResponse.getData();
            }
        } catch (Exception e) {
            logger.error("get pvc failed ",  e.getMessage());
        }
```

这样就可以了.
