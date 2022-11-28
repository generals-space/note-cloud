# kuber-Ingress对象中的Nginx配置

参考文章

1. [413 error with Kubernetes and Nginx ingress controller](https://stackoverflow.com/questions/49918313/413-error-with-kubernetes-and-nginx-ingress-controller)
2. [官方文档 NGINX Configuration](https://github.com/kubernetes/ingress-nginx/blob/master/docs/user-guide/nginx-configuration/index.md)
3. [nginx-configuration](https://github.com/kubernetes/ingress-nginx/tree/master/docs/user-guide/nginx-configuration)

通过`docker push`向`harbor`推送镜像时出现了413错误, 于是想起还需要定义`client_max_body_size`字段解除这个限制.

```
$ d push harbor.generals.space/base/centos7:1.0
The push refers to repository [harbor.generals.space/base/centos7]
163a76a9e552: Layer already exists
36906321f3c4: Pushing [==================================================>]  185.1MB/185.1MB
f9911f13d28d: Layer already exists
1d31b5806ba4: Pushing [==================================================>]  199.8MB
error parsing HTTP 413 response body: invalid character '<' looking for beginning of value: "<html>\r\n<head><title>413 Request Entity Too Large</title></head>\r\n<body>\r\n<center><h1>413 Request Entity Too Large</h1></center>\r\n<hr><center>openresty/1.15.8.1</center>\r\n</body>\r\n</html>\r\n"
```

那么问题来了, ingress里应该在哪里改nginx的配置?

按照参考文章1的采纳答案, 可以在Ingress配置文件中添加注解

```yaml
apiVersion: extensions/v1beta1
kind: Ingress
metadata:
  name: harbor-ing
  namespace: harbor
  annotations:
    nginx.ingress.kubernetes.io/proxy-body-size: "0"
```

经实验有效.

这是把nginx的字段重新命名了啊, 如果想配置其他字段怎么办? 

按照参考文章2, 一共有3种方法可以修改nginx配置: 

1. config map用于修改全局配置, 这个config map在部署ingress时就在deploy中通过`--configmap`选项引入了, 所以直接修改config map然后apply, 再重启pod就可以生效.
2. 就是上面的annotation, 可以修改单个ingress对象的配置.
3. 没仔细看...

至于像nginx字段到ingress字段的名称的对应关系, 可以到官方文档查找, 这里就不详细介绍了.
