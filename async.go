package pool

// AsyncJob 定义异步作业
type AsyncJob struct {
	*SyncJob
}

// NewAsyncJob 新建异步作业
func NewAsyncJob(id string, level Level) *AsyncJob {
	return &AsyncJob{
		SyncJob: NewSyncJob(id, level),
	}
}

// Wait 异步作业无需等待
func (*AsyncJob) Wait() interface{} {
	panic("implement me")
}
