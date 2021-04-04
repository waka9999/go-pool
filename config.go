package pool

import "fmt"

const (
	queueCapacityDefault = 200
	queueCapacityMin     = 100
	queueCapacityMax     = 500

	jobCapacityDefault = 50
	jobCapacityMin     = 20
	jobCapacityMax     = 100

	workerCapacityDefault = 20
	workerCapacityMin     = 10
	workerCapacityMax     = 100

	queueTimeout = 3000
)

// Config 调度池设置项
type Config struct {
	// queue 容量
	QueueCapacity int `yaml:"queueCapacity"`

	// job 容量
	JobCapacity int `yaml:"jobCapacity"`

	// worker 容量
	WorkerCapacity int `yaml:"workerCapacity"`

	// timeout 排队超时
	QueueTimeout int `yaml:"queueTimeout"`
}

// DefaultConfig 调度池默认设置
func DefaultConfig() *Config {
	return &Config{
		QueueCapacity:  queueCapacityDefault,
		JobCapacity:    jobCapacityDefault,
		WorkerCapacity: workerCapacityDefault,
		QueueTimeout:   queueTimeout,
	}
}

// Check 设置项检查
func (c *Config) Check() {
	if c.QueueCapacity < queueCapacityMin {
		c.QueueCapacity = queueCapacityMin
	}
	if c.QueueCapacity > queueCapacityMax {
		c.QueueCapacity = queueCapacityMax
	}

	if c.JobCapacity < jobCapacityMin {
		c.JobCapacity = jobCapacityMin
	}
	if c.JobCapacity > jobCapacityMax {
		c.JobCapacity = jobCapacityMax
	}

	if c.WorkerCapacity < workerCapacityMin {
		c.WorkerCapacity = workerCapacityMin
	}
	if c.WorkerCapacity > workerCapacityMax {
		c.WorkerCapacity = workerCapacityMax
	}
}

// 重写 String 方法
func (c *Config) String() string {
	return fmt.Sprintf(
		"\nPool Config:"+
			"\n\tQueueCapacity: %d"+
			"\n\tJobCapacity： %d"+
			"\n\tWorkerCapacity: %d"+
			"\n\tQueueTimeout: %d\n",
		c.QueueCapacity,
		c.JobCapacity,
		c.WorkerCapacity,
		c.QueueTimeout)
}
