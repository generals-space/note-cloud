参考文章

1. [kubernetes-client/java patch示例](https://github.com/kubernetes-client/java/blob/client-java-parent-9.0.0/examples/src/main/java/io/kubernetes/client/examples/PatchExample.java)
    - PATCH_FORMAT_JSON_PATCH
    - PATCH_FORMAT_STRATEGIC_MERGE_PATCH
    - PATCH_FORMAT_APPLY_YAML
2. [k8s-client(java)从6.0.1升级到11.0.0出现patch问题may not be specified for non-apply patch/cannot unmarshal...](https://blog.csdn.net/qq_33999844/article/details/115279872)
    - 一般的jsonPatch（标准规则）是有5种，add/remove/move/replace/copy，k8s是不支持move/copy的
3. [11.0.0 patchNamespacedCustomObject exception “PatchOptions.meta.k8s.io...”](https://github.com/kubernetes-client/java/issues/1575)

java client: 9.0.0

kube: v1.16.2

对主机`label`的新增操作.

```java
Map<String, Object> labelOperator = new HashMap<>();
labelOperator.put("op", "add");
labelOperator.put("path", "/metadata/labels/key01");
labelOperator.put("value", "val01");

List<Object> body = new ArrayList<>();
body.add(labelOperator);

CoreV1Api api = null;
try {
    PatchUtils.patch(
        Object.class,
        () -> api.patchNodeCall(
            nodeName, newV1Patch(jsonUtils.objToStr(body)),
            null, null, null, null, null
        ),
        V1Patch.PATCH_FORMAT_JSON_PATCH,
        api.getApiClient()
    );
    return Boolean.TRUE
}catch (ApiException e) {

    return Boolean.FALSE
}
```
