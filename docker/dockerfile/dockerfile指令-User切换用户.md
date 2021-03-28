# dockerfile指令-User切换用户

参考文章

1. [USER 指定当前用户](https://www.cntofu.com/book/139/image/dockerfile/user.md)

dockerfile中的`User`指令只能用于切换用户(要求目标用户事先存在), 之后的`RUN`指令都会以该用户执行, 即, 只在dockerfile构建过程中有效.

ta不能创建用户, 也不能在 CMD 中使用, 如果希望在容器启动时切换用户执行, 最好使用`su`, `runuser`等指令.
