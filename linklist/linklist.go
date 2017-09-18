package linklist
import "errors"

type Linklist struct {
	head *node	
	tail *node
}
var EmptyLinkList = errors.New("the linklist is empty")
func (l *Linklist) InsertHead(ele interface {}) (error) {
	node := newNode(ele)
	currentFirstNode := l.head.next
	if currentFirstNode == nil {
		l.head.next = node
		node.pre = l.head
		l.tail = node
		return nil
	}
	
	currentFirstNode.pre = node
	node.next = currentFirstNode
	node.pre =  l.head
	l.head.next = node
	return nil
}

func (l *Linklist) InsertTail(ele interface {}) (error) {
	node := newNode(ele)
	if l.tail == nil {
		l.head.next = node
		node.pre = l.head
		l.tail = node
		return nil
	}
	
	l.tail.next = node
	node.pre = l.tail
	l.tail = node
	return nil
}

func (l *Linklist) GetHead() (interface {}, error) {
	if l.IsEmpty() {
		return "", EmptyLinkList
	} 
	return l.head.next.ele, nil
}

func (l *Linklist) GetTail() (interface {}, error) {
	if l.IsEmpty() {
		return "", EmptyLinkList
	}
	return l.tail.ele, nil
}

func (l *Linklist) ExtractHead() (interface {}, error) {
	if l.IsEmpty() {
		return "", EmptyLinkList
	}
	ele := l.head.next.ele
	if l.head.next == l.tail {
		l.head.next = nil
		l.tail = nil
		return ele, nil
	}
	nextNode := l.head.next.next
	l.head.next = nextNode
	nextNode.pre = l.head
	return ele, nil
}

func (l *Linklist) ExtractTail() (interface {}, error) {
	if l.IsEmpty() {
		return "", EmptyLinkList
	}
	ele := l.tail.ele
	if l.head.next == l.tail{
		l.head.next = nil
		l.tail = nil
		return ele, nil
	}
	preNode := l.tail.pre
	l.tail = preNode
	l.tail.next = nil
	return ele, nil
}

// func (l *Linklist) GetHeadNode() node {
// 	return l.head
// }

func (l *Linklist) IsEmpty() bool {
	if l.head.next == nil {
		return true
	}
	return false
}
func NewLinklist() (*Linklist) {
	return &Linklist {
		head : newNode("gakki"),
		tail : nil,
	}
} 