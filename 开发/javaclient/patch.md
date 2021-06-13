参考文章

1. [kubernetes-client/java patch示例](https://github.com/kubernetes-client/java/blob/client-java-parent-9.0.0/examples/src/main/java/io/kubernetes/client/examples/PatchExample.java)
    - PATCH_FORMAT_JSON_PATCH
    - PATCH_FORMAT_STRATEGIC_MERGE_PATCH
    - PATCH_FORMAT_APPLY_YAML

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
