# crd开发(二) go module工程

参考文章

1. [go modules support](https://github.com/kubernetes/code-generator/issues/57)
2. [Generator only works if k8s.io/apimachinery is in the $GOPATH](https://github.com/kubernetes/code-generator/issues/21)
1. [code-generator使用](https://tangxusc.github.io/blog/2019/05/code-generator%E4%BD%BF%E7%94%A8/)
    - `code-generator`中各工具的作用: `deepcopy-gen`, `client-gen`, `informer-gen`, `lister-gen`
    - 实战向, 示例步骤极为详细, 可作参考(不过ta的示例代码中各文件路径与其他文章有些区别, 可进行比较后采用)
    - 使用`go mod`创建`CRD`工程

前面在尝试CRD开发的时候, 发现`code-generator`生成的代码, 目标路径总是在`GOPATH`目录下, 并且没有办法写绝对路径.

```log
$ bash vendor/k8s.io/code-generator/generate-groups.sh all $(pwd)/pkg/client $(pwd)/pkg/apis mycrdgroup:v1
Generating deepcopy funcs
F1224 04:39:46.533441  125875 main.go:82] Error: Failed making a parser: unable to add directory "/home/project/mycrd/pkg/apis/mycrdgroup/v1": unable to import "/home/project/mycrd/pkg/apis/mycrdgroup/v1": import "/home/project/mycrd/pkg/apis/mycrdgroup/v1": cannot import absolute path
```

这就导致了我们只能先把工程放在`GOPATH`下, 代码生成完成后, 再移出来进行`go mod`初始化与依赖管理.

但是, 如果后期我们需要修改CRD内容, 重新生成代码时呢? 

我在网上找了很久, 没找到有对这个话题有明确说明的文章, 最后终于找到了参考文章3. 

其实准备步骤基本没有差别, 只是在生成代码时需要注意以下几点:

1. 可以直接使用`go module`模式创建工程, 并用`go mod init 工程名`初始化工程;
2. 使用`go mod`创建的工程, 就不能再使用`$GOPATH/src/k8s.io/code-generator/generate-groups.sh`来调用脚本了, 我们需要把`code-generator`工程添加到`go.mod`文件中, 但不是直接修改此文件, 而是在`hack/tools.go`文件中引入`code-generator`包, 相当于一个占位符, 可以查看一个`sample-controller`工程中此文件的内容, 还是十分容易理解的.
3. 最重要是, `generate-groups.sh`脚本有一个选项`--output-base`, 用于指定生成代码的路径. 此选项的详细介绍可以见`code-generator/hack/update-codegen.sh`. `update-codegen.sh`可以借鉴`sample-controller`工程的内容, 执行路径是与此相关的.

几个额外的知识点:

`code-generator/hack/boilerplate.go.txt`和`custom-boilerplate.go.txt`, 这两个文件中是一段将写到自动生成代码中的内容, 一般是`License`的说明. 可以通过`generate-groups.sh`脚本的`--go-header-file`指定路径.

