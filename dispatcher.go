// pool 包实现 job&worker 线程池
package pool

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

type dispatcherState = uint64

const (
	// 停止状态
	stopped dispatcherState = iota

	// 开始状态
	started
)

// Dispatcher 定义调度器
type Dispatcher struct {
	// 调度器设置
	config *Config

	// 调度器状态
	state dispatcherState

	// 作业队列
	jobQueue *jobQueue

	// 工作者池
	workerPool chan *worker

	// 队列超时
	timeout time.Duration

	// 退出通道
	exitCh chan struct{}

	// 互斥锁
	mu sync.Mutex

	waitGroup sync.WaitGroup
}

// NewDispatcher 新建调度器
func NewDispatcher(config *Config) *Dispatcher {
	return &Dispatcher{
		config:  config,
		state:   stopped,
		timeout: time.Duration(config.QueueTimeout) * time.Millisecond,
		exitCh:  make(chan struct{}),
	}
}

// Start 调度器启动
func (d *Dispatcher) Start() {
	d.mu.Lock()
	defer d.mu.Unlock()

	// 避免重复启动, 转换状态为开始状态
	if atomic.CompareAndSwapUint64(&d.state, stopped, started) {
		// 初始化退出通道
		d.exitCh = make(chan struct{})

		// 新建作业队列
		d.jobQueue = newJobQueue(d.config.QueueCapacity, d.config.JobCapacity, d.exitCh)

		// 新建工作者池
		d.workerPool = make(chan *worker, d.config.WorkerCapacity)
		for i := 0; i < d.config.WorkerCapacity; i++ {
			d.workerPool <- newWorker(fmt.Sprintf("worker-%d", i), d.exitCh)
		}

		// 初始化 JobQueue
		d.waitGroup.Add(1)
		go func() {
			d.jobQueue.run()
			d.waitGroup.Done()
		}()

		d.waitGroup.Add(1)
		// 调度器运行
		go func() {
			d.run()
			d.waitGroup.Done()
		}()

	}
}

// run 调度器运行
func (d *Dispatcher) run() {
	for {
		select {
		case <-d.exitCh:
			return
		case job := <-d.jobQueue.jobsCh:

			if job == nil {
				continue
			}

			ctx, cancel := context.WithTimeout(context.Background(), d.timeout)
			defer cancel()
			go func() {
				select {
				case <-ctx.Done():
					if ctx.Err() == context.DeadlineExceeded {
						job.Error(ctx.Err())
					}
				case worker := <-d.workerPool:
					if atomic.LoadUint64(&d.state) == started {
						if ctx.Err() == context.DeadlineExceeded {
							job.Error(ctx.Err())
							d.workerPool <- worker
						} else {
							worker.start(job)
							d.workerPool <- worker
						}
					}
				}
			}()
		}
	}
}

// Join 加入作业
func (d *Dispatcher) Join(j Job) {
	d.mu.Lock()
	defer d.mu.Unlock()

	// 调度器处于开始状态
	if atomic.LoadUint64(&d.state) == started {
		select {
		case <-d.exitCh:
			return
		default:
			d.jobQueue.insert(j)
			return
		}
	}
}

// Stop 调度器停止
func (d *Dispatcher) Stop() {
	// 调度器转换为停止状态
	if atomic.CompareAndSwapUint64(&d.state, started, stopped) {
		// 关闭退出通道
		close(d.exitCh)

		// 等待
		d.waitGroup.Wait()

		// 调度器队列清除
		d.jobQueue.clear()
	}
}
