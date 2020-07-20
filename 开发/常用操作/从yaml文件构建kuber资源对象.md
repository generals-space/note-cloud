# 从yaml文件构建kuber资源对象

参考文章

1. [How to deserialize Kubernetes YAML file](https://stackoverflow.com/questions/44306554/how-to-deserialize-kubernetes-yaml-file)
    - 问题和采纳答案都太旧了, `Ramiro Berrelleza`的回答才是本文的目标.

kuber版本: 1.16.2

`go.mod`文件内容如下

```go
module gosts

go 1.13

require (
	github.com/elazarl/goproxy v0.0.0-20180725130230-947c36da3153 // indirect
	github.com/gogo/protobuf v1.3.1 // indirect
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b
	github.com/google/gofuzz v1.1.0 // indirect
	github.com/googleapis/gnostic v0.1.0 // indirect
	github.com/imdario/mergo v0.3.9 // indirect
	github.com/onsi/ginkgo v1.11.0 // indirect
	golang.org/x/oauth2 v0.0.0-20200107190931-bf48bf16ab8d // indirect
	golang.org/x/sys v0.0.0-20191022100944-742c48ecaeb7 // indirect
	golang.org/x/time v0.0.0-20200630173020-3af7569d3a1e // indirect
	k8s.io/api v0.17.0
	k8s.io/apimachinery v0.17.0
	k8s.io/client-go v0.17.0
	k8s.io/kube-openapi v0.0.0-20200410145947-61e04a5be9a6 // indirect
	k8s.io/utils v0.0.0-20200619165400-6e3d28b6ed19 // indirect
	sigs.k8s.io/structured-merge-diff/v3 v3.0.0 // indirect
	sigs.k8s.io/yaml v1.2.0 // indirect
)
```

```go
package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/golang/glog"
	apiAppsv1 "k8s.io/api/apps/v1"
	apimYaml "k8s.io/apimachinery/pkg/util/yaml"
	cgKuber "k8s.io/client-go/kubernetes"
	cgClientcmd "k8s.io/client-go/tools/clientcmd"
	cgHomedir "k8s.io/client-go/util/homedir"
)

func main() {
	// 先尝试从 ~/.kube 目录下获取配置, 如果没有, 则尝试寻找 Pod 内置的认证配置
	var kubeconfig string
	home := cgHomedir.HomeDir()
	kubeconfig = filepath.Join(home, ".kube", "config")
	cfg, err := cgClientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		glog.Infof("failed to build kubeconfig: %s", err.Error())
		return
	}

	kuberClient, err := cgKuber.NewForConfig(cfg)
	if err != nil {
		glog.Infof("failed to get kube client: %s", err.Error())
		return
	}

	sts := &apiAppsv1.StatefulSet{}
	yamlbytes, err := ioutil.ReadFile("sts.yaml")
	if err != nil {
		glog.Infof("failed to read sts file: %s", err.Error())
		return
	}

	reader := bytes.NewReader(yamlbytes)
	decoder := apimYaml.NewYAMLOrJSONDecoder(reader, len(yamlbytes))
	err = decoder.Decode(sts)
	if err != nil {
		glog.Infof("failed to parse sts file: %s", err.Error())
		return
	}
	fmt.Printf("%+v\n", sts)

	stsObj, err := kuberClient.AppsV1().StatefulSets("default").Create(sts)
	if err != nil {
		glog.Infof("failed to create sts: %s", err.Error())
		return
	}
	fmt.Printf("=== sts: %+v\n", stsObj)
}
```
