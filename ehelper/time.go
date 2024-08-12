package ehelper

import (
	"math/rand"
	"time"
)

func RandomSleep(minMills, maxMills int) {
	r := Int(minMills, maxMills)
	time.Sleep(time.Duration(r) * time.Millisecond)
}

func Int(min int, max int) int {
	rand.Seed(time.Now().UnixNano())
	//rand.New(rand.NewSource(time.Now().UnixNano()))
	time.Sleep(10 * time.Nanosecond)
	return min + rand.Intn(max-min)
}
