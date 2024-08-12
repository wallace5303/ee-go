package etask

import (
	"context"
	"reflect"
	"sync"
	"time"

	"github.com/wallace5303/ee-go/ehelper"
	"github.com/wallace5303/ee-go/elog"
	"github.com/wallace5303/ee-go/eruntime"
	"github.com/wallace5303/ee-go/eutil"
)

// uniqueActions 描述了唯一的任务，即队列中只能存在一个在执行的任务。
var UniqueActions = []string{}

var (
	taskQueue         []*Task
	queueLock         = sync.Mutex{}
	currentTaskAction string
)

type Task struct {
	Action  string
	Handler reflect.Value
	Args    []interface{}
	Created time.Time
	Timeout time.Duration
}

// 添加任务
// action: "task.demo"
// uniqueAction: "unique.task.demo"
func AddTask(action string, handler interface{}, args ...interface{}) {
	AddTaskWithTimeout(action, handler, 24*time.Hour, args...)
}

func AddTaskWithTimeout(action string, handler interface{}, timeout time.Duration, args ...interface{}) {
	if eruntime.IsExiting {
		return
	}

	// 是否为唯一任务
	currentActions := getCurrentActions()
	if ehelper.Contains(action, currentActions) && ehelper.Contains(action, UniqueActions) {
		return
	}

	queueLock.Lock()
	defer queueLock.Unlock()
	taskQueue = append(taskQueue, &Task{
		Action:  action,
		Timeout: timeout,
		Handler: reflect.ValueOf(handler),
		Args:    args,
		Created: time.Now(),
	})
}

// 获取当前队列中所有任务action
func getCurrentActions() (ret []string) {
	queueLock.Lock()
	defer queueLock.Unlock()

	if currentTaskAction != "" {
		ret = append(ret, currentTaskAction)
	}

	for _, task := range taskQueue {
		ret = append(ret, task.Action)
	}
	return
}

func Contain(action string, moreActions ...string) bool {
	actions := append(moreActions, action)
	actions = ehelper.RemoveDuplicatedElem(actions)

	for _, task := range taskQueue {
		if ehelper.Contains(task.Action, actions) {
			return true
		}
	}
	return false
}

func Status() {
	tasks := taskQueue
	var items []map[string]any
	count := map[string]int{}
	for _, task := range tasks {
		action := task.Action
		count[action]++

		item := map[string]any{"action": action}
		items = append(items, item)
	}

	// 添加正在执行的任务
	if currentTaskAction != "" {
		items = append([]map[string]any{{"action": currentTaskAction}}, items...)
	}

	if len(items) < 1 {
		items = []map[string]any{}
	}

	if eruntime.IsDev() {
		elog.Logger.Infoln("task status data: ", items)
	}
}

func ExecTask() {
	task := popTask()
	if task == nil {
		return
	}

	if eruntime.IsExiting {
		return
	}

	execTask(task)
}

func popTask() (ret *Task) {
	queueLock.Lock()
	defer queueLock.Unlock()

	if len(taskQueue) == 0 {
		return
	}

	ret = taskQueue[0]
	taskQueue = taskQueue[1:]
	return
}

func execTask(task *Task) {
	defer eutil.Recover()

	args := make([]reflect.Value, len(task.Args))
	for i, v := range task.Args {
		if v == nil {
			args[i] = reflect.New(task.Handler.Type().In(i)).Elem()
		} else {
			args[i] = reflect.ValueOf(v)
		}
	}

	currentTaskAction = task.Action

	ctx, cancel := context.WithTimeout(context.Background(), task.Timeout)
	defer cancel()
	ch := make(chan bool, 1)
	go func() {
		task.Handler.Call(args)
		ch <- true
	}()

	select {
	case <-ctx.Done():
		elog.Logger.Warnf("task [%s] timeout", task.Action)
	case <-ch:
		elog.Logger.Infof("task [%s] done", task.Action)
	}

	currentTaskAction = ""
}
