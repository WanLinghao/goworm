package main

import "goworm/visitor"

func main() {
	ch := make(chan byte, 2)
	visitor.BeginWork(ch)
	<-ch
	//ch is to keep main goroutine waitting for subgoroutine
} 