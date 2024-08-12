package eapp

import (
	"os"
	"os/signal"
	"reflect"
	"sync"
	"syscall"
	"time"

	"github.com/wallace5303/ee-go/ehelper"
	"github.com/wallace5303/ee-go/elog"
	"github.com/wallace5303/ee-go/eruntime"
	"github.com/wallace5303/ee-go/eutil"
)

var (
	exitLock        = sync.Mutex{}
	eventsMap       = make(map[string]*Staff)
	BeforeCloseFlag = "beforeClose"
)

var Events = []string{
	BeforeCloseFlag,
}

type Staff struct {
	Name    string
	Handler reflect.Value
	Args    []interface{}
}

func Register(eventName string, handler interface{}, args ...interface{}) {
	if !ehelper.SlicesContains(Events, eventName) {
		elog.CoreLogger.Warnf("[ee-go] event: %s is not supported", eventName)
		return
	}

	eventsMap[eventName] = &Staff{
		Name:    eventName,
		Handler: reflect.ValueOf(handler),
		Args:    args,
	}
}

func execStaff(staff *Staff) {
	defer eutil.Recover()

	args := make([]reflect.Value, len(staff.Args))
	for i, v := range staff.Args {
		if v == nil {
			args[i] = reflect.New(staff.Handler.Type().In(i)).Elem()
		} else {
			args[i] = reflect.ValueOf(v)
		}
	}

	staff.Handler.Call(args)
}

// Run the program and listen for Signal
func Run() {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	sig := <-sigCh

	elog.CoreLogger.Infof("[ee-go] received signal: %s", sig)
	Close()
}

// Close process
func Close() {
	exitLock.Lock()
	defer func() {
		exitLock.Unlock()
	}()

	eruntime.IsExiting = true
	elog.CoreLogger.Infof("[ee-go] process is exiting...")

	// [todo] wait other
	if eventsMap[BeforeCloseFlag] != nil {
		execStaff(eventsMap[BeforeCloseFlag])
	}

	go func() {
		time.Sleep(500 * time.Millisecond)
		eruntime.IsExiting = false
		elog.CoreLogger.Infof("[ee-go] process has exited!")
		os.Exit(0)
	}()
}
