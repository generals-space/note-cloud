使用 kubebuilder, 貌似无需直接通过 client-go 与 apiserver 进行交互.

另外又由于封装了一层, 貌似也不再需要 WorkQueue 操作了...

不再有 Informer 机制.

