# indexer索引器

参考文章

1. [client-go系列之4---Indexer](https://zhuanlan.zhihu.com/p/266512431)
    - 其中 2.1.1 节的索引器数据结构值得一看.
2. [client-go Indexer索引器](https://herbguo.gitbook.io/client-go/informer#4.2-indexer-suo-yin-qi)
    - index 的使用示例...有点难懂
3. [client-go 之 Indexer 的理解](https://blog.51cto.com/u_15077560/2584555)
    - 这篇文章的 indexer 示例代码比参考文章2的易懂.

```
// 包含的所有索引器/分类以及对应的实现
Indexers: {  
    "namespace": NamespaceIndexFunc,
    "nodeName": NodeNameIndexFunc,
}

// 包含的所有索引分类中所有的索引数据
Indices: {
    //namespace 这个索引分类下的所有索引数据
    "namespace": {  
         // Index 就是一个索引键下所有的对象键列表
        "default": ["pod-1", "pod-2"], 
        // Index
        "kube-system": ["pod-3"]   
    },

    //nodeName 这个索引分类下的所有索引数据(对象键列表)
    "nodeName": {
         // Index
        "node1": ["pod-1"],
         // Index
        "node2": ["pod-2", "pod-3"]
    }
}
```
