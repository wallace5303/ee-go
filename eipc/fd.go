// go:build !windows
package eipc

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"syscall"

	"github.com/wallace5303/ee-go/ehelper"
)

const (
	// node.js与子进程通信的文件描述符
	NODE_CHANNEL_FD = "NODE_CHANNEL_FD"

	// 固定的通信的消息类型
	listeningEvent  = "listening"
	connectionEvent = "connection"
	closeEvent      = "close"
	errorEvent      = "error"
)

var ListenEvents = []string{
	listeningEvent,
	connectionEvent,
	closeEvent,
	errorEvent,
}

// 通信桥
type Bridge struct {
	fdReader  *FDReader
	bufReader *bufio.Reader
	events    map[string]EventHandler
	closeChan chan int
	connected bool
}

// file descriptor reader
type FDReader struct {
	fileD *os.File
}

// 实现 io.Reader  的 Read 方法
func (br *FDReader) Read(p []byte) (n int, err error) {
	length := syscall.CmsgSpace(4)
	mbuf := make([]byte, length)
	n, _, _, _, err = syscall.Recvmsg(int(br.fileD.Fd()), p, mbuf, 0)
	if err != nil {
		return 0, err
	}
	return n, nil
}

// 通信的上下文
type Context struct {
	Id        string
	Data      string
	bridgeRef *Bridge
}

// 通信的消息
type Message struct {
	Id   string `json:"id"`
	Type string `json:"type"`
	Data string `json:"data"`
}

// 事件函数
type EventHandler func(Context)

// 与 node.js 建立连接
func Connect() *Bridge {
	nodeChannelFD := os.Getenv(NODE_CHANNEL_FD)
	nodeChannelFDInt, _ := strconv.Atoi(nodeChannelFD)
	fileD := os.NewFile(uintptr(int(nodeChannelFDInt)), "ipc"+nodeChannelFD)
	ior := &FDReader{fileD: fileD}

	b := &Bridge{
		fdReader:  ior,
		bufReader: bufio.NewReader(ior),
		events:    make(map[string]EventHandler),
	}
	return b
}

// 注册事件函数
func (bridge *Bridge) On(event string, fn EventHandler) {
	bridge.events[event] = fn
}

func (bridge *Bridge) SendByType(data string, msgType string) error {
	fd := int(bridge.fdReader.fileD.Fd())
	responseMsg := Message{
		Id:   ehelper.GetRandomString(5),
		Data: data,
		Type: msgType,
	}
	jsonData, _ := json.Marshal(responseMsg)
	return syscall.Sendmsg(fd, append(jsonData, '\n'), nil, nil, 0)
}

func (ctx *Context) Response(data string) error {
	fd := int(ctx.bridgeRef.fdReader.fileD.Fd())
	responseMsg := Message{
		Id:   ctx.Id,
		Data: data,
		Type: "response",
	}
	jsonData, _ := json.Marshal(responseMsg)
	return syscall.Sendmsg(fd, append(jsonData, '\n'), nil, nil, 0)
}

// 监听
func (bridge *Bridge) Listen() {
	bridge.closeChan = make(chan int)
	go func() {
		bridge.listen()
	}()
	bridge.SendByType("", connectionEvent)
	<-bridge.closeChan
}

func (bridge *Bridge) listen() {
	for {
		bridge.readMessage()
	}
}

func (bridge *Bridge) readMessage() (int, error) {
	data, err := bridge.bufReader.ReadBytes('\n')
	if err != nil {
		return 0, err
	}
	msg := new(Message)
	json.Unmarshal(data, msg)

	fmt.Println("msg:", msg)
	eventName := msg.Type

	if eventName == "ready" {
		eventName = connectionEvent
	} else if eventName == "close" {
		bridge.closeChan <- 1
	} else if !bridge.connected {
		return 0, nil
	}

	eventListener, exists := bridge.events[eventName]
	if !exists {
		return 0, nil
	}
	go eventListener(Context{
		Id:        msg.Id,
		Data:      msg.Data,
		bridgeRef: bridge,
	})
	if eventName == connectionEvent {
		bridge.connected = true
	}
	return 1, nil
}
