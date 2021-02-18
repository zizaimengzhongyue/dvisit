# dvisit
支持通过字符串路径对于 golang 数据进行动态读写操作：
```
package main

import (
	"fmt"

	"github.com/zizaimengzhongyue/dvisit"
)

type Test struct {
	Key   string
	Slice []int
}

func main() {
	test := Test{Key: "key", Slice: []int{1, 2, 3}}

	key, err := dvisit.Get(test, "Key")
	if err != nil {
		panic(err)
	}
	fmt.Println(key)

	val, err := dvisit.Get(test, "Slice.0")
	if err != nil {
		panic(err)
	}
	fmt.Println(val)

    err = dvisit.Set(&test, "Key", "new key")
    if err != nil {
        panic(err)
    }
	key, err = dvisit.Get(test, "Key")
	if err != nil {
		panic(err)
	}
	fmt.Println(key)
}
```

暂时不支持通配符  
