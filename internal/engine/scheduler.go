package engine

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/golang/glog"
	"github.com/livepeer/swarm-chaos/internal/model"
)

type (
	interval struct {
		min time.Duration
		max time.Duration
	}

	filter struct {
		key   string
		value string
	}

	task struct {
		interval  interval
		operation model.OperationType
		filter    filter
		running   bool
	}

	// Scheduler executes tasks toward playground
	Scheduler struct {
		playground model.Playground
		running    bool
		context    context.Context
		cancel     context.CancelFunc
		tasks      []task
	}
)

// NewScheduler creates a new Scheduler
func NewScheduler(playground model.Playground) *Scheduler {
	// p := make([]model.Playground, len(playgrounds))
	// copy(p, playgrounds)
	// return &Scheduler{playgrounds: p}
	return &Scheduler{playground: playground}
}

// AddPlayground adds a new playground
// func (sc *Scheduler) AddPlayground(playground model.Playground) {
// 	sc.playgrounds = append(sc.playgrounds, playground)
// }

// ScheduleTask ...
func (sc *Scheduler) ScheduleTask(intervalFrom, intervalTo string, operation model.OperationType, filterKey, filterValue string) error {
	var interval interval
	pd, err := time.ParseDuration(intervalFrom)
	if err != nil {
		return err
	}
	interval.min = pd
	pd, err = time.ParseDuration(intervalTo)
	if err != nil {
		return err
	}
	interval.max = pd
	task := task{
		interval: interval,
		filter: filter{
			key:   filterKey,
			value: filterValue,
		},
		operation: operation,
	}
	sc.tasks = append(sc.tasks, task)
	return nil
}

// ClearTasks stops all tasks and clears tasks list
func (sc *Scheduler) ClearTasks() error {
	sc.StopTasks()
	sc.tasks = make([]task, 0)
	return nil
}

// StartTasks starts scheduled tasks
func (sc *Scheduler) StartTasks() error {
	if sc.running {
		return fmt.Errorf("Already started")
	}
	ctx, cancel := context.WithCancel(context.Background())
	sc.context = ctx
	sc.cancel = cancel
	for _, task := range sc.tasks {
		go sc.startTaskLoop(ctx, &task)
	}
	sc.running = true
	glog.Infof("Started %d tasks", len(sc.tasks))

	return nil
}

func (sc *Scheduler) entitiesByLabel(key, value string) ([]model.Entity, error) {
	res := make([]model.Entity, 0)
	entities, err := sc.playground.Entities()
	if err != nil {
		return nil, err
	}
	for _, e := range entities {
		labels := e.Labels()
		if labels[key] == value {
			res = append(res, e)
		}
	}
	return res, nil
}

func (sc *Scheduler) startTaskLoop(ctx context.Context, task *task) {
	for {
		toWait := task.interval.min + time.Duration(rand.Int63n(int64(task.interval.max-task.interval.min)))
		glog.Infof("Wating %s", toWait)
		time.Sleep(toWait)
		select {
		case <-ctx.Done():
			return
		default:
		}
		glog.Infof("Finding entities with label %s:%s", task.filter.key, task.filter.value)
		entities, err := sc.entitiesByLabel(task.filter.key, task.filter.value)
		if err != nil {
			glog.Infof("Can't get entities: %v", err)
			return
		}
		glog.Infof("Found %d entities", len(entities))
		if len(entities) == 0 {
			glog.Infof("No entities found")
			continue
		}
		ent := entities[rand.Intn(len(entities))]
		glog.Infof("Killing entity %s", ent.Name())
		err = ent.Do(task.operation)
		if err != nil {
			glog.Infof("Error killing entity %s: %v", ent.Name(), err)
		}
	}
}

// StopTasks stops the scheduler
func (sc *Scheduler) StopTasks() bool {
	if sc.cancel != nil {
		sc.cancel()
		sc.running = false
		sc.cancel = nil
		sc.context = nil
		return true
	}
	return false
}
