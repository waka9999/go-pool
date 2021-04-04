package pool

// Level 定义作业级别
type Level = int64

const (
	// High 优先级
	High Level = iota

	// App 应用级
	App

	// Task 任务级
	Task

	// System 系统级
	System

	// Low 低优先级
	Low
)

// Job 定义作业接口
type Job interface {
	// 作业 ID
	ID() string

	// 作业级别
	Level() Level

	// 作业超时执行
	Error(err error)

	// 正常作业执行
	Do()

	// 同步等待
	Wait() interface{}
}
