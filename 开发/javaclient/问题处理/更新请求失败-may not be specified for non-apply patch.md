# 更新请求失败-may not be specified for non-apply patch

参考文章

1. [k8s-client(java)从6.0.1升级到11.0.0出现patch问题may not be specified for non-apply patch/cannot unmarshal...](https://blog.csdn.net/qq_33999844/article/details/115279872)

```java
import io.kubernetes.client.custom.V1Patch;
import io.kubernetes.client.openapi.JSON;
import io.kubernetes.client.util.PatchUtils;
```
