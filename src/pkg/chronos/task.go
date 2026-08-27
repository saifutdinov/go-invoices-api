package chronos

import (
	"context"
	"time"
)

type ScheduledTask interface {
	GetStartTime() time.Time
	Run(context.Context)
}

type FuncTask struct {
	StartTime time.Time
	Handler   func(context.Context)
}

func NewTask(
	startTime time.Time,
	handler func(context.Context),
) ScheduledTask {
	return &FuncTask{
		StartTime: startTime,
		Handler:   handler,
	}
}

func (t *FuncTask) GetStartTime() time.Time {
	return t.StartTime
}

func (t *FuncTask) Run(ctx context.Context) {
	if t.Handler == nil {
		return
	}

	t.Handler(ctx)
}

type ScheduledTaskQueue []ScheduledTask

func (stq ScheduledTaskQueue) Len() int {
	return len(stq)
}

func (stq ScheduledTaskQueue) Less(i, j int) bool {
	return stq[i].GetStartTime().Before(
		stq[j].GetStartTime(),
	)
}

func (stq ScheduledTaskQueue) Swap(i, j int) {
	stq[i], stq[j] = stq[j], stq[i]
}

func (stq *ScheduledTaskQueue) Push(x any) {
	*stq = append(
		*stq,
		x.(ScheduledTask),
	)
}

func (stq *ScheduledTaskQueue) Pop() any {
	queue := *stq

	n := len(queue)

	task := queue[n-1]

	queue[n-1] = nil

	*stq = queue[:n-1]

	return task
}

func (stq *ScheduledTaskQueue) Top() ScheduledTask {
	if stq.Len() == 0 {
		return nil
	}

	return (*stq)[0]
}
