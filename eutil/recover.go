package eutil

import (
	"bytes"
	"fmt"
	"os"
	"runtime"

	"github.com/wallace5303/ee-go/elog"
)

var (
	dunno     = []byte("???")
	centerDot = []byte("·")
	dot       = []byte(".")
	slash     = []byte("/")
)

func Recover() {
	if r := recover(); r != nil {
		stack := stack()
		msg := fmt.Sprintf("PANIC RECOVERED: %v\n\t%s\n", r, stack)
		elog.Logger.Errorf(msg)
	}
}

// stack
// 可以借鉴errors包的打印堆栈实现
func stack() []byte {
	buf := &bytes.Buffer{}
	var lines [][]byte
	var lastFile string
	for i := 2; ; i++ {
		pc, file, line, ok := runtime.Caller(i)
		// fmt.Printf("--- %s:%d (0x%x) ok:%t, \n", file, line, pc, ok)
		if !ok {
			break
		}
		fmt.Fprintf(buf, "%s:%d (0x%x)\n", file, line, pc)
		if file != lastFile {
			data, err := os.ReadFile(file)
			if err != nil {
				continue
			}
			lines = bytes.Split(data, []byte{'\n'})
			lastFile = file
		}
		line--
		fmt.Fprintf(buf, "\t%s: %s\n", function(pc), source(lines, line))
	}
	return buf.Bytes()
}

// function
func function(pc uintptr) []byte {
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return dunno
	}
	name := []byte(fn.Name())
	if lastslash := bytes.LastIndex(name, slash); lastslash >= 0 {
		name = name[lastslash+1:]
	}
	if period := bytes.Index(name, dot); period >= 0 {
		name = name[period+1:]
	}
	name = bytes.Replace(name, centerDot, dot, -1)
	return name
}

// source
func source(lines [][]byte, n int) []byte {
	if n < 0 || n >= len(lines) {
		return dunno
	}
	return bytes.Trim(lines[n], " \t")
}
