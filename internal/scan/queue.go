package scan

type Node struct {
	next  *Node
	token *Token
}

type Queue struct {
	front *Node
	end   *Node
	size  uint64
}

func NewQueue() *Queue {
	queue := new(Queue)
	return queue
}

func (q *Queue) Add(token *Token) {
	if q.front == nil {
		node := new(Node)
		node.token = token
		q.front = node
		q.end = node
	} else {
		node := new(Node)
		node.token = token
		q.end.next = node
		q.end = q.end.next
	}
	q.size++
}

func (q *Queue) Clear() {
	q.front = nil
	q.end = nil
	q.size = 0
}

func (q *Queue) RemoveFront() *Token {
	if q.front == nil {
		return nil
	} else {
		node := q.front.token
		q.front = q.front.next
		q.size--
		return node
	}
}

func (q *Queue) Len() uint64 {
	return q.size
}

func (q *Queue) Peak() *Token {
	return q.front.token
}

func (q *Queue) String() string {
	s := ""
	for tmp := q.front; tmp != nil; tmp = tmp.next {
		s += tmp.token.String() + " "
	}
	return s
}
