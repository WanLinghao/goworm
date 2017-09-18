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



var liveNessCond = sync.NewCond(&sync.Mutex{})
var processorDone = false
var visitorDone = false
var urlQueue = queue.NewQueue()
var contentQueue = queue.NewQueue()
var visitedPool = &sync.Map{}
var errorPage = &sync.Map{}

func end(ch chan byte) {
	ch <- 1
}
func Visit(ch chan byte) error {
	defer end(ch)
	for ;true; {
	//	fmt.Printf("visitor\n")
		if urlQueue.IsEmpty() {
				liveNessCond.L.Lock()
				if processorDone {
					//processor has finished and watting 
					if contentQueue.IsEmpty() && urlQueue.IsEmpty() {
						visitorDone = true
						liveNessCond.L.Unlock()
						liveNessCond.Broadcast()
						
						return nil
					} 
					//liveNessCond.L.Unlock()
					liveNessCond.Broadcast()
					visitorDone = true
					fmt.Printf("waitting here!\n")
					//liveNessCond.L.Unlock()
					liveNessCond.Wait()//wait here
					liveNessCond.L.Unlock()
					fmt.Printf("end watting\n")
					if urlQueue.IsEmpty() {
						liveNessCond.Broadcast()
						fmt.Printf("return here\n")
						return nil
					}
				} else {
					//processor is working
					visitorDone = true
					
					liveNessCond.Wait()
					liveNessCond.L.Unlock()
					if processorDone {
						if urlQueue.IsEmpty() {
						
							return nil
						} else {
							visitorDone = false
							
							//liveNessCond.Broadcast()
						}
					}
				}
			
		}
		//fmt.Printf("gakki\n")
 		noneTypeEle, err := urlQueue.Dequeue()
		if err != nil {
			fmt.Printf("error !!!!!!!\n")
			return err
		}

		joint, ok := ConvertToUrlJoint(noneTypeEle)
		if !ok{
			return errors.New("can't transfer to urlJoint type")
		}
		
		if !ShouldVisit(joint) {
			continue
		}
		
		visitedPool.Store(joint.containURL, 0)
		fmt.Printf("url is %s\n",joint.containURL)
		resp, err := http.Get(joint.containURL)
		
		if err != nil {
			fmt.Printf("error happens:%s\n", err)
			continue
		}
		if resp.StatusCode == 404  {
			//this is the first time we visit this page, so errorPage must have no this page key
			writeToFile(joint.contentURL, joint.containURL, "/home/wlh/gakkiki")
			fmt.Printf("error 404: %s\n", joint.containURL)
			tmp := &sync.Map{}
			tmp.Store(joint.contentURL, 0)
			errorPage.Store(joint.containURL, tmp)
			continue
		} else if resp.StatusCode != 200 {
			fmt.Printf("error %d: %s\n",resp.StatusCode, joint.containURL)
		}


		body, err := ioutil.ReadAll(resp.Body)
		str := string(body)
	//	fmt.Printf("%s\n",str)
		// if contentQueue.IsEmpty(){
		// 	contentQueue.Enqueue(contentJoint{URL:joint.containURL, content:str,})
		// 	liveNessCond.Broadcast()

		// }
		contentQueue.Enqueue(contentJoint{URL:joint.containURL, content:str,})
		//resp.Body.Close()
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
				fmt.Printf("get lockggg\n")
				liveNessCond.L.Lock()
				fmt.Printf("get lock\n")
				if visitorDone {
					//visitor has finished and watting 
					if contentQueue.IsEmpty() && urlQueue.IsEmpty() {
						processorDone = true
						liveNessCond.L.Unlock()
						liveNessCond.Broadcast()
						return nil
					} 
					//liveNessCond.L.Unlock()
					fmt.Printf("processor watting here!\n")
					liveNessCond.Broadcast()
					processorDone = true
					liveNessCond.Wait()
					liveNessCond.L.Unlock()
					if contentQueue.IsEmpty() {
						liveNessCond.Broadcast()
						return nil
					}
				} else {
					//visitor is working
					fmt.Printf("processor 1\n")
					processorDone = true
					liveNessCond.Broadcast()
					fmt.Printf("wati at 1\n")
					liveNessCond.Wait()
					liveNessCond.L.Unlock()
					fmt.Printf("wait end 1\n")
					if visitorDone {
						if contentQueue.IsEmpty() {
							return nil
						} else {
							fmt.Printf("enter 2\n")
							processorDone = false
							//liveNessCond.L.Unlock()
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

		// rule := "'https?://kubernetes.io[^\"']*?'"
		// reg := regexp.MustCompile(rule)
		// slice := reg.FindAllString(content.content, -1)
		// Doo(content.URL,slice)
		
		// rule = "\"https?://kubernetes.io[^\"']*?\""
		// reg = regexp.MustCompile(rule)
		// slice = reg.FindAllString(content.content, -1)
		// Doo(content.URL,slice)

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

func Doo(str string,slice []string) {
	for _, s := range slice {
			fmt.Printf("%s\n", s)
			if strings.HasPrefix(s, "\"") || strings.HasPrefix(s, "'") {
				s = s[1:]
			}
			if strings.HasSuffix(s, "\"") || strings.HasSuffix(s, "'") {
				s = s[:len(s)-1]
			}
			urlQueue.Enqueue(urlJoint{
				contentURL : str,
				containURL : s,
			})
		}
}

func BeginWork(ch chan byte) bool {
	urlQueue.Enqueue(urlJoint{
		contentURL : "gakki",
		containURL : "https://kubernetes.io/",
	})
	
	go Visit(ch)
	go Process(ch)
	
	return true
	// fmt.Printf("%s", visitErr)
	// fmt.Printf("%s", processErr)
}

func writeToFile(rawPageURL, invalidURL, filePath string )  {
	fd, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		fmt.Printf("write file error!\n")
		return 
	}
	

	line1 := []byte("rawPageURL: "+rawPageURL+"\n")
	fd.Write(line1)
	line2 := []byte("invalidURL: "+invalidURL+"\n")
	fd.Write(line2)
	if err := fd.Close(); err != nil {
		fmt.Printf("close file error!\n")
	}

}


