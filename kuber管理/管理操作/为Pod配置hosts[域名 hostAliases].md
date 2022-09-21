# 为Pod配置hosts[域名]

参考文章

1. [k8s给pod添加hosts](https://www.cnblogs.com/route/p/16119323.html)

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
## ...省略
spec:
  template:
    metadata:
    ## ...省略
    spec:
      hostAliases:
      - ip: "10.236.9.220"
        hostnames:
        - "xxx.xxx.local"
      containers:
      - name: 
        image: 

```
