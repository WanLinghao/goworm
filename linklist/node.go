package linklist
type node struct {
	ele interface {}
	next *node
	pre *node	
}

func newNode(ele interface {}) *node{
	return &node {
		ele: ele,
		next: nil,
		pre: nil,
	}
}
