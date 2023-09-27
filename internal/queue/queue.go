package queue

import (
	"errors"
	"sync"
)

var (
	ErrQueueIsFull  = errors.New("Очередь полная")
	//ErrQueueisEmpty = errors.New("Очередь пустая")
)

type Queue struct {
	c    *sync.Cond
	cap  int
	buff []interface{}
}

func New(cap int) *Queue {
	return &Queue{
		cap:  cap,
		buff: make([]interface{}, 0, cap),
		c:    sync.NewCond(&sync.Mutex{}),
	}
}

// Publish добавляет в очередь, если она полная возращает ошибку.
func (q *Queue) Publish(msg interface{}) error {
	q.c.L.Lock()
	defer q.c.L.Unlock()
	if len(q.buff) < q.cap {
		q.buff = append(q.buff, msg)
		q.c.Signal()
		return nil
	}

	return ErrQueueIsFull
}
// Consume читает из очереди, если очередь пустая, то блокируется пока не сможет считать.
func (q *Queue) Consume() interface{} {
	q.c.L.Lock()
	defer q.c.L.Unlock()
	if len(q.buff) <= 0 {
		q.c.Wait()
	}
	msg := q.buff[0]
	q.buff = q.buff[1:]

	return msg
}

func (q *Queue) Export() []interface{} {
	return q.buff
}
