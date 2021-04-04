package pool

import (
	"container/heap"
	"sync"

	queue "github.com/waka9999/go-queue"
)

// jobQueue 定义作业队列
type jobQueue struct {
	// job 通道
	jobsCh chan Job

	// 优先级队列
	queue *queue.PriorityQueue

	// 队列互斥锁
	mu sync.Mutex

	// ready 通道
	readyCh chan struct{}

	// 退出通道
	exitCh chan struct{}
}

// newJobQueue 新建作业队列
func newJobQueue(qCap int, jCap int, ech chan struct{}) *jobQueue {
	return &jobQueue{
		jobsCh:  make(chan Job, jCap),
		queue:   queue.NewPriorityQueue(qCap),
		readyCh: make(chan struct{}, jCap),
		exitCh:  ech,
	}
}

// run 获取队列中的首个作业放到 ready 通道
func (jq *jobQueue) run() {
	for {
		select {
		case <-jq.exitCh:
			return
		case <-jq.readyCh:
			jq.mu.Lock()
			jq.jobsCh <- jq.queue.First().Value().(Job)
			heap.Pop(jq.queue)
			jq.mu.Unlock()
		}
	}
}

// insert 插入作业
func (jq *jobQueue) insert(j Job) {
	select {
	case <-jq.exitCh:
		return
	default:
		// 新建队列存储项
		item := queue.NewItem(j, j.Level())

		jq.mu.Lock()
		heap.Push(jq.queue, item)
		jq.mu.Unlock()

		jq.readyCh <- struct{}{}
	}
}

// 清除作业队列
func (jq *jobQueue) clear() {
	jq.mu.Lock()
	defer jq.mu.Unlock()
	close(jq.jobsCh)
	close(jq.readyCh)
	jq.queue.Clear()
}
