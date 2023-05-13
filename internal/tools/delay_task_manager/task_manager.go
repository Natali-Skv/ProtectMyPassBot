package delay_task_manager

import (
	"context"
	"time"
)

type Processed interface {
	Process()
}

type DelayTaskManager interface {
	AddTask(processed Processed)
	Run(ctx context.Context, period time.Duration)
}
