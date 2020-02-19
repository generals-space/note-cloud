调用`code-generator`生成代码, 直接执行`./hack/update-codegen.sh`即可, 执行路径没有要求.

另外, 最好在建立工程之初就创建`go.mod`文件, 安装`code-generator`依赖(在`hack/tools.go`文件中作为占位引用). 因为`main.go`和`controller.go`中引用了`code-generator`生成的代码, 如果在创建这两个文件后, 且在生成代码前, 使用`go mod download`, 会因为这两个文件的依赖包不存在而失败.
