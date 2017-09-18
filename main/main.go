package main

//import "fmt"
//import gakki "worm/Goworm/queue"
//import "fmt"
//import "worm/Goworm/linklist"
//import "worm/Goworm/queue"
import "worm/Goworm/visitor"
//import "sync"
 //import "net/http"
// import "io/ioutil"
 //import "regexp"

//var test = sync.Map{}
func main() {
	
	ch := make(chan byte, 2)
	visitor.BeginWork(ch)
	<-ch
	

	// resp, err := http.Get("https://kubernetes.io/")
	// fmt.Printf("resp: %#v", resp)
	// fmt.Printf("err: %s",err)



} 