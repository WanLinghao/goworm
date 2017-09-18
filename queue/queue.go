package queue

import "goworm/linklist"
import "sync"

type Queue struct {
	list *linklist.Linklist
	lock *sync.Mutex
}

func (q *Queue) Dequeue() (interface{}, error) {
	q.lock.Lock()
	ele, err := q.list.ExtractTail()
	q.lock.Unlock()
	return ele, err
}

func (q *Queue) Enqueue(ele interface{}) error {
	q.lock.Lock()
	q.list.InsertHead(ele)
	q.lock.Unlock()
	return nil
}

func (q *Queue) IsEmpty() bool {
	q.lock.Lock()
	ele := q.list.IsEmpty()
	q.lock.Unlock()
	return ele
}

func NewQueue() *Queue {
	return &Queue{
		list: linklist.NewLinklist(),
		lock: &sync.Mutex{},
	}
}
