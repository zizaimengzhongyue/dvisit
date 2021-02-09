# dvisit
支持通过字符串路径对于 golang 数据进行动态访问：
```
package main

import (
    "fmt"
)

type Test struct {
    Key string
    Slice []int
}

test := Test{Key: "key", Slice: []int{1,2,3}}

key, err := dvisit.Get("key")
if err != nil {
    panic(err)
}
fmt.Println(key)

val, err := dvisit.Get("Slice.0")
if err != nil {
    panic(err)
}
fmt.Println(val)
```

暂时不支持通配符  
