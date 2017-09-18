package main

import "goworm/visitor"

func main() {
	filePath := "/home/wlh/gakkiki"
	ch := make(chan byte, 2)
	visitor.BeginWork(ch, filePath)
	<-ch
	//ch is to keep main goroutine waitting for subgoroutine
}
