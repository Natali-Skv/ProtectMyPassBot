package delay_task_manager

import (
	"context"
	"github.com/Natali-Skv/ProtectMyPassBot/internal/tools/delay_task_manager"
	"go.uber.org/zap"
	"time"
)

type DelayTaskManager struct {
	l     *zap.Logger
	Tasks []Timed
}

type Timed struct {
	AddedTs time.Time
	delay_task_manager.Processed
}

func NewDelayTaskManager(l *zap.Logger) *DelayTaskManager {
	return &DelayTaskManager{l: l, Tasks: make([]Timed, 0)}
}

func (tm *DelayTaskManager) AddTask(processed delay_task_manager.Processed) {
	tm.Tasks = append(tm.Tasks, Timed{AddedTs: time.Now(), Processed: processed})
}

func (tm *DelayTaskManager) Run(ctx context.Context, period time.Duration) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if len(tm.Tasks) > 0 && time.Now().After(tm.Tasks[0].AddedTs.Add(period)) {
				tm.Tasks[0].Process()
				tm.Tasks = tm.Tasks[1:]
			}
		}
	}
}
