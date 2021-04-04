package pool

// SyncJob 定义同步作业
type SyncJob struct {
	id    string
	level Level
}

// NewSyncJob 新建同步作业
func NewSyncJob(id string, level Level) *SyncJob {
	return &SyncJob{
		id:    id,
		level: level,
	}
}

// ID 获取作业ID
func (j *SyncJob) ID() string {
	return j.id
}

// Level 获取作业级别
func (j *SyncJob) Level() Level {
	return j.level
}
