参考文章

1. [Installing Trivy](https://aquasecurity.github.io/trivy/v0.38/getting-started/installation/)

trivy 可以扫描多种场景下的漏洞, 包括容器镜像, 整个k8s集群, 及本地文件系统等.

```log
$ trivy

Usage:
  trivy [global flags] command [flags] target
  trivy [command]

Examples:
  # Scan a container image
  $ trivy image python:3.4-alpine

  # Scan local filesystem
  $ trivy fs .

Available Commands:
  filesystem  Scan local filesystem
  image       Scan a container image
  kubernetes  [EXPERIMENTAL] Scan kubernetes cluster
  vm          [EXPERIMENTAL] Scan a virtual machine image
```

## 扫描镜像

```log
$ trivy image docker.io/calico/cni:v3.24.1 -q

opt/cni/bin/host-local (gobinary) ## 容器中的二进制文件

Total: 1 (UNKNOWN: 0, LOW: 0, MEDIUM: 1, HIGH: 0, CRITICAL: 0)

┌──────────────────┬────────────────┬──────────┬───────────────────┬───────────────┬───────┐
│     Library      │ Vulnerability  │ Severity │ Installed Version │ Fixed Version │ Title │
├──────────────────┼────────────────┼──────────┼───────────────────┼───────────────┼───────┤
│ golang.org/x/sys │ CVE-2022-29526 │ MEDIUM   │      v0.0.0       │     0.0.0     │       │
└──────────────────┴────────────────┴──────────┴───────────────────┴───────────────┴───────┘

opt/cni/bin/install (gobinary) ## 容器中的二进制文件

Total: 5 (UNKNOWN: 0, LOW: 0, MEDIUM: 1, HIGH: 3, CRITICAL: 1)

┌────────────────────────────────┬────────────────┬──────────┬───────────────────┬───────────────┬────────┐
│            Library             │ Vulnerability  │ Severity │ Installed Version │ Fixed Version │  Title │
├────────────────────────────────┼────────────────┼──────────┼───────────────────┼───────────────┼────────┤
│ github.com/emicklei/go-restful │ CVE-2022-1996  │ CRITICAL │      v2.11.2      │ 2.16.0        │        │
├────────────────────────────────┼────────────────┼──────────┼───────────────────┼───────────────┼────────┤
│ golang.org/x/net               │ CVE-2022-27664 │ HIGH     │      v0.0.0       │ 0.0.0-20220906│ golang │
│                                ├────────────────┤          │                   ├───────────────┼────────┤
│                                │ CVE-2022-41723 │          │                   │ 0.7.0         │ golang │
│                                ├────────────────┼──────────┤                   ├───────────────┼────────┤
│                                │ CVE-2022-41717 │ MEDIUM   │                   │ 0.4.0         │ golang │
├────────────────────────────────┼────────────────┼──────────┼───────────────────┼───────────────┼────────┤
│ golang.org/x/text              │ CVE-2022-32149 │ HIGH     │      v0.3.7       │ 0.3.8         │ golang │
└────────────────────────────────┴────────────────┴──────────┴───────────────────┴───────────────┴────────┘
```

ta扫描了此镜像中所有可执行文件, 并分别给出了ta们拥有的漏洞.

## `--severity`只输出高危漏洞

上面扫描出的结果太多了, 有时候我们只想知道目标镜像是否为高危, 对于那种"中等"威胁的镜像, 可以直接忽略.

可以使用`-s`选项进行过滤.

```
-s, --severity string        severities of security issues to be displayed (comma separated) (default "UNKNOWN,LOW,MEDIUM,HIGH,CRITICAL")
```

```log
$ trivy image docker.io/calico/cni:v3.24.1 -q -s 'HIGH,CRITICAL'

opt/cni/bin/install (gobinary) ## 容器中的二进制文件

Total: 4 (HIGH: 3, CRITICAL: 1)

┌────────────────────────────────┬────────────────┬──────────┬───────────────────┬───────────────┬────────┐
│            Library             │ Vulnerability  │ Severity │ Installed Version │ Fixed Version │  Title │
├────────────────────────────────┼────────────────┼──────────┼───────────────────┼───────────────┼────────┤
│ github.com/emicklei/go-restful │ CVE-2022-1996  │ CRITICAL │      v2.11.2      │ 2.16.0        │        │
├────────────────────────────────┼────────────────┼──────────┼───────────────────┼───────────────┼────────┤
│ golang.org/x/net               │ CVE-2022-27664 │ HIGH     │      v0.0.0       │ 0.0.0-20220906│ golang │
│                                ├────────────────┤          │                   ├───────────────┼────────┤
│                                │ CVE-2022-41723 │          │                   │ 0.7.0         │ golang │
├────────────────────────────────┼────────────────┼──────────┼───────────────────┼───────────────┼────────┤
│ golang.org/x/text              │ CVE-2022-32149 │ HIGH     │      v0.3.7       │ 0.3.8         │ golang │
└────────────────────────────────┴────────────────┴──────────┴───────────────────┴───────────────┴────────┘
```

