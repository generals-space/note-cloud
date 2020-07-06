# docker update更新资源限制

`--memory`: 可使用 10m, 1g这种格式. 在设置内存时, 不可比当前的 swap 值小, 如果一定要缩小的话, 需要再加上`--memory-swap`选项, 用于指定目标 swap 的大小. 但这个选项并不单指 swap, 而是目标 memory+swap 的和.

