package etask

import (
	"time"

	"github.com/wallace5303/ee-go/ehelper"
	"github.com/wallace5303/ee-go/eutil"
)

// 每隔一段时间执行一次
func Every(interval time.Duration, f func()) {
	ehelper.RandomSleep(50, 200)
	for {
		func() {
			defer eutil.Recover()
			f()
		}()

		time.Sleep(interval)
	}
}
