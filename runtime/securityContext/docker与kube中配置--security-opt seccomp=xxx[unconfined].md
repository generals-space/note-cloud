参考文章

1. [Docker run reference](https://docs.docker.com/engine/reference/run/#security-configuration)
2. [Seccomp security profiles for Docker](https://docs.docker.com/engine/security/seccomp/)

## docker

```
docker run --security-opt seccomp=unconfined -e MYSQL_ROOT_PASSWORD=123456 mysql:latest
```

## kube

要先创建一个`SeccompProfile`

```yaml
apiVersion: seccomp.security.alpha.kubernetes.io/v1
kind: SeccompProfile
metadata:
  name: unconfined
spec:
  seccompProfile:
    type: Unconfined
```

然后在 pod 中指定

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: my-pod
spec:
  containers:
  - name: my-container
    image: my-image
    securityContext:
      seccompProfile:
        type: Localhost
        localhostProfile: unconfined
```

注意: 这里的 type 属性为 `Localhost`, 表示使用的是本地 seccompProfile. 如果要使用集群中的 seccompProfile, 可以将 type 属性设置为 RuntimeDefault.

