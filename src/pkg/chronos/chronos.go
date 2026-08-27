package chronos

import (
	"container/heap"
	"context"
	"sync"
	"time"
)

const defaultQueueSize = 100

type Chronos struct {
	taskQueue   *ScheduledTaskQueue
	taskChannel chan ScheduledTask
	stopChannel chan struct{}

	eventAlarm *time.Timer

	stopOnce sync.Once
}

func New(ctx context.Context) *Chronos {
	queue := new(ScheduledTaskQueue)
	heap.Init(queue)

	c := &Chronos{
		taskQueue:   queue,
		taskChannel: make(chan ScheduledTask, defaultQueueSize),
		stopChannel: make(chan struct{}),
		eventAlarm:  time.NewTimer(time.Hour),
	}

	go c.wait(ctx)

	return c
}

func (c *Chronos) Schedule(task ScheduledTask) {
	select {
	case c.taskChannel <- task:
	case <-c.stopChannel:
	}
}

func (c *Chronos) Stop() {
	c.stopOnce.Do(func() {
		close(c.stopChannel)
	})
}

func (c *Chronos) pushTask(task ScheduledTask) {
	heap.Push(c.taskQueue, task)
}

func (c *Chronos) popTask() ScheduledTask {
	return heap.Pop(c.taskQueue).(ScheduledTask)
}

func (c *Chronos) resetAlarm() {
	if !c.eventAlarm.Stop() {
		select {
		case <-c.eventAlarm.C:
		default:
		}
	}

	if c.taskQueue.Len() == 0 {
		c.eventAlarm.Reset(time.Hour)
		return
	}

	delay := time.Until(
		c.taskQueue.Top().GetStartTime(),
	)

	if delay < 0 {
		delay = 0
	}

	c.eventAlarm.Reset(delay)
}

func (c *Chronos) wait(ctx context.Context) {
	for {
		c.resetAlarm()

		select {
		case <-ctx.Done():
			return

		case <-c.stopChannel:
			return

		case task := <-c.taskChannel:
			if task == nil {
				continue
			}

			c.pushTask(task)

		case <-c.eventAlarm.C:
			if c.taskQueue.Len() == 0 {
				continue
			}

			if time.Now().Before(
				c.taskQueue.Top().GetStartTime(),
			) {
				continue
			}

			task := c.popTask()

			go func() {
				task.Run(ctx)
			}()
		}
	}
}
