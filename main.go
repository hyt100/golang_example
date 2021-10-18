package main

import (
	"example/persistent" // 注意：自定义包的路径写法是按照 “相对于GOPATH/src” 的路径写的
	"fmt"
)

func main() {
	fmt.Println("Run Test: ")
	persistent.Test()
}
