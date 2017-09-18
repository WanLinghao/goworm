package visitor

import (
	"net/http"
	"fmt"
	"io/ioutil"
	"sync"
	"regexp"
	"worm/Goworm/queue"
	"errors"
	"strings"
	"os"

)
//liveNessCond is a cond to cooperate between processor and visitor
var liveNessCond = sync.NewCond(&sync.Mutex{})

//processorDone is a flag to indicate if the processor has finished its work and been waitting
//so be visitorDone
var processorDone = false
var visitorDone = false


var urlQueue = queue.NewQueue()
var contentQueue = queue.NewQueue()
var visitedPool = &sync.Map{}
var errorPage = &sync.Map{}

func end(ch chan byte) {
	ch <- 1
}
func Visit(ch chan byte, fd *os.File) error {
	defer end(ch)
	for ;true; {
		if urlQueue.IsEmpty() {
				liveNessCond.L.Lock()
				if processorDone {
					//processor is watting  
					if contentQueue.IsEmpty() && urlQueue.IsEmpty() {
						visitorDone = true
						liveNessCond.L.Unlock()
						liveNessCond.Broadcast()
						return nil
					} 

					liveNessCond.Broadcast()
					visitorDone = true
					liveNessCond.Wait()
					liveNessCond.L.Unlock()

					if urlQueue.IsEmpty() {
						liveNessCond.Broadcast()
						return nil
					}
				} else {
					//processor is working
					visitorDone = true
				//	liveNessCond.Broadcast()
					liveNessCond.Wait()
					liveNessCond.L.Unlock()
					if processorDone {
						if urlQueue.IsEmpty() {
							//liveNessCond.Broadcast()
							return nil
						} else {
							visitorDone = false
							//liveNessCond.Broadcast()						
						}
					}
				}
		}
 		noneTypeEle, err := urlQueue.Dequeue()
		if err != nil {
			fmt.Printf("dequeue error: %s\n", err)
			return err
		}

		joint, ok := ConvertToUrlJoint(noneTypeEle)
		if !ok{
			return errors.New("can't transfer to urlJoint type\n")
		}
		
		if !ShouldVisit(joint) {
			continue
		}
		
		visitedPool.Store(joint.containURL, 0)
		fmt.Printf("url is %s\n",joint.containURL)
		resp, err := http.Get(joint.containURL)
		if err != nil {
			fmt.Printf("visit url error: %s\n", err)
			continue
		}
		if resp.StatusCode == 404  {
			//this is the first time we visit this page, so errorPage must have no this page key
			writeToFile(joint.contentURL, joint.containURL, fd)
			fmt.Printf("catch 404 error: %s\n", joint.containURL)
			tmp := &sync.Map{}
			tmp.Store(joint.contentURL, 0)
			errorPage.Store(joint.containURL, tmp)
			continue
		} else if resp.StatusCode != 200 {
			fmt.Printf("catch %d error: %s\n",resp.StatusCode, joint.containURL)
		}


		body, err := ioutil.ReadAll(resp.Body)
		str := string(body)
		contentQueue.Enqueue(contentJoint{URL:joint.containURL, content:str,})

	}
	return nil
}

func ShouldVisit(joint urlJoint) bool {
	if _, ok := visitedPool.Load(joint.containURL); !ok{
		//this url hasn't been visited
		return true
	}
	emptyTypeEle, ok := errorPage.Load(joint.containURL)
	if ok {
		//this url has been visited and it is an error page
		ele, ok2 := emptyTypeEle.(sync.Map)
		if ok2{
			ele.Store(joint.contentURL, 0)
		} else {
			fmt.Printf("error,can't tarnsfer to sync.Map")
		}
	}	
	return false 
}

func Process (ch chan byte) error {
	for ;true; {
		if contentQueue.IsEmpty() {
				liveNessCond.L.Lock()
				if visitorDone {
					//visitor has finished and watting 
					if contentQueue.IsEmpty() && urlQueue.IsEmpty() {
						//both queue is empty now and the processor has finished its job ,return
						processorDone = true
						liveNessCond.L.Unlock()
						liveNessCond.Broadcast()
						return nil
					} 

					liveNessCond.Broadcast()
					processorDone = true
					liveNessCond.Wait()
					liveNessCond.L.Unlock()

					if contentQueue.IsEmpty() {
						liveNessCond.Broadcast()
						return nil
					}
				} else {
					//visitor is working, 
					processorDone = true
					liveNessCond.Broadcast()
					liveNessCond.Wait()
					liveNessCond.L.Unlock()
					if visitorDone {
						if contentQueue.IsEmpty() {
							//in case visitor is waitting
							//liveNessCond.Broadcast()
							return nil
						} else {
							processorDone = false
							//liveNessCond.Broadcast()
						}
					}
				}
			
		}
		
		emptyTypeEle, err:= contentQueue.Dequeue()
		if err != nil {
			return err
		}
		content, ok := ConvertToContentJoint(emptyTypeEle)
		if !ok {
			return errors.New("error,can't conver to type ContentJoint")
		}

		rule := "href=\"[^\"']*\""
		reg := regexp.MustCompile(rule)
		slice := reg.FindAllString(content.content, -1)
		for _, s := range slice {
			s = strings.TrimPrefix(s, "href=\"")
			s = strings.TrimSuffix(s, "\"")
			if strings.HasPrefix(s, "https://kubernetes.io") {
				urlQueue.Enqueue(urlJoint{
				contentURL : content.URL,
				containURL : s,
			})
			} else if strings.HasPrefix(s, "/"){
				urlQueue.Enqueue(urlJoint{
				contentURL : content.URL,
				containURL : "https://kubernetes.io" + s,
			})
			} else {
				continue
			}
		}
	}
	end(ch)
	return nil
}

func BeginWork(ch chan byte) bool {
	urlQueue.Enqueue(urlJoint{
		contentURL : "gakki",
		containURL : "https://kubernetes.io/",
	})
	
	filePath := "/home/wlh/gakkiki"
	fd, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		fmt.Printf("error: %s\n", err)
		return false
	}

	go Visit(ch, fd)
	go Process(ch)

	if err := fd.Close(); err != nil {
		fmt.Printf("close file error!\n")
	}
	return true
}

func writeToFile(rawPageURL, invalidURL string, fd *os.File)  {
	line1 := []byte("rawPageURL: "+rawPageURL+"\n")
	fd.Write(line1)
	line2 := []byte("invalidURL: "+invalidURL+"\n")
	fd.Write(line2)
}


