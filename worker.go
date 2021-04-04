package pool

// worker 定义工作者
type worker struct {
	// worker id
	id string

	// 退出通道
	exitCh <-chan struct{}
}

// newWorker 新建工作者
func newWorker(id string, ech <-chan struct{}) *worker {
	return &worker{
		id:     id,
		exitCh: ech,
	}
}

// start 工作者开始作业
func (w *worker) start(job Job) {
	select {
	case <-w.exitCh:
		return
	default:
		job.Do()
	}
}
