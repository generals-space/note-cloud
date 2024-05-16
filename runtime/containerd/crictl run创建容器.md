参考文章

1. [CRICTL User Guide](https://github.com/containerd/containerd/blob/main/docs/cri/crictl.md)
    - `/etc/crictl.yaml`配置

crictl 创建一个容器, 比 docker 麻烦了很多啊...所有容器都都运行在 sandbox 中.

需要先创建一个 sandbox.

```json
cat <<EOF | tee sandbox.json
{
    "metadata": {
        "name": "nginx-sandbox",
        "namespace": "default"
    }
}
EOF
```

```
crictl runp --runtime runsc sandbox.json
```
