# kuber-在Pod command字段中引用环境变量[env]

参考文章

1. [为容器设置启动时要执行的命令和参数](https://kubernetes.io/zh/docs/tasks/inject-data-application/define-command-argument-container/)
2. [Kubernetes(k8s)为容器设置启动时要执行的命令和参数](https://www.orchome.com/9877)

在`command`中引用环境变量时, 需要使用`$()`符号, 而非`$`或`${}`.

```yaml
    spec:
      containers:
      - name: centos7
        env:
        - name: FILE
          value: /etc/os-release
        command: ["tail", "-f", "$(FILE)"]
```

进入到容器内部, 使用`ps`查看进程列表, 会发现进程的启动命令中, 环境变量已经被替换掉了.

```log
sh-4.2# ps -ef
UID         PID   PPID  C STIME TTY          TIME CMD
root          1      0  0 15:37 ?        00:00:00 tail -f /etc/os-release
```

------

在阅读 kubernetes v1.16.0 源码时, 正好有这么一段函数.

```go
func ExpandContainerCommandAndArgs(
	container *v1.Container, envs []EnvVar,
) (command []string, args []string) {
	mapping := expansion.MappingFuncFor(EnvVarsToMap(envs))

	if len(container.Command) != 0 {
		for _, cmd := range container.Command {
			command = append(command, expansion.Expand(cmd, mapping))
		}
	}

	if len(container.Args) != 0 {
		for _, arg := range container.Args {
			args = append(args, expansion.Expand(arg, mapping))
		}
	}

	return command, args
}
func (m *kubeGenericRuntimeManager) generateContainerConfig(
	container *v1.Container, 
	pod *v1.Pod, 
	restartCount int, 
	podIP, 
	imageRef string,
) (*runtimeapi.ContainerConfig, func(), error) {
    ...
	command, args := kubecontainer.ExpandContainerCommandAndArgs(container, opts.Envs)
}
```
