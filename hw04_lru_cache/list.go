package hw04lrucache

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	head   *ListItem
	tail   *ListItem
	length int
}

func (l *list) Len() int {
	return l.length
}

func (l *list) Front() *ListItem {
	return l.head
}

func (l *list) Back() *ListItem {
	return l.tail
}

func (l *list) PushFront(v interface{}) *ListItem {
	prevHead := l.head
	l.head = &ListItem{Value: v, Next: prevHead}

	if prevHead != nil {
		prevHead.Prev = l.head
	}

	if l.length == 0 {
		l.tail = l.head
	}

	l.length++

	return l.head
}

func (l *list) PushBack(v interface{}) *ListItem {
	prevTail := l.tail
	l.tail = &ListItem{Value: v, Prev: prevTail}

	if prevTail != nil {
		prevTail.Next = l.tail
	}

	if l.length == 0 {
		l.head = l.tail
	}

	l.length++

	return l.tail
}

func (l *list) Remove(item *ListItem) {
	prev := item.Prev
	next := item.Next

	if prev == nil {
		l.head = next
	} else {
		prev.Next = next
	}

	if next == nil {
		l.tail = prev
	} else {
		next.Prev = prev
	}

	l.length--
}

func (l *list) MoveToFront(item *ListItem) {
	prev := item.Prev

	if prev == nil {
		return
	}

	next := item.Next
	prev.Next = next

	if next == nil {
		l.tail = prev
	} else {
		next.Prev = prev
	}

	prevHead := l.head

	item.Prev = nil
	item.Next = prevHead

	l.head = item
	prevHead.Prev = item
}

func NewList() List {
	return new(list)
}
