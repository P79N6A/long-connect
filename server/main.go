package main

import (
	"fmt"
	"im/comet"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	fmt.Printf("CPU NUM IS: %d\n", runtime.NumCPU())
	// 启动长连接清理线程
	go comet.ClearConnList()
	go comet.CurrentLongConnection()
	comet.StartServer("tcp", ":9999")
}
