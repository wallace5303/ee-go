package eipc

import (
	"bufio"
	"encoding/json"
	"os"
	"strconv"
	"syscall"

	"github.com/wallace5303/ee-go/elog"
)

// node.js childProcess.spawn() ipc
const (
	NODE_CHANNEL_FD = "NODE_CHANNEL_FD"
)

type Bridge struct {
	lbr         *BridgeReader
	eventMap    map[string]EventCallback
	closeChan   chan int
	established bool
}

type BridgeReader struct {
	fd     *os.File
	reader *bufio.Reader
}

type Context struct {
	Id      string
	Data    string
	Options map[string]string
	lbRef   *Bridge
}

type EventCallback func(Context)

type Message struct {
	Id      string `json:"id"`
	MsgType string `json:"type"`
	Data    string `json:"data"`
}

func Init() {
	elog.CoreLogger.Infof("[ee-go] load ipc")

}

// implement io.Reader Read
func (lbr *BridgeReader) Read(p []byte) (n int, err error) {
	// 后去传递的数据的有效负载长度
	length := syscall.CmsgSpace(4)
	mbuf := make([]byte, length)
	// 将收到的数据从内核空间拷贝至用户空间
	// oob: Out of band data https://beej.us/298C/oob_overview.html
	n, _, _, _, err = syscall.Recvmsg(int(lbr.fd.Fd()), p, mbuf, 0)
	if err != nil {
		return 0, err
	}
	return n, nil
}

func Establish() *Bridge {
	reader := establishChannel()
	lb := &Bridge{
		lbr:      reader,
		eventMap: make(map[string]EventCallback),
	}
	return lb
}

func establishChannel() *BridgeReader {
	nodeChannelFD := os.Getenv(NODE_CHANNEL_FD)
	nodeChannelFDInt, _ := strconv.Atoi(nodeChannelFD)
	fd := os.NewFile(uintptr(int(nodeChannelFDInt)), "lbipc"+nodeChannelFD)
	reader := &BridgeReader{fd: fd}
	reader.reader = bufio.NewReader(reader)
	return reader
}

func (lb *Bridge) On(event string, callback EventCallback) {
	lb.eventMap[event] = callback
}

func (lb *Bridge) Listen() {
	lb.closeChan = make(chan int)
	go func() {
		lb.listen()
	}()
	lb.sendByType("ready", "establish")
	<-lb.closeChan
}

func (lb *Bridge) sendByType(data string, msgType string) error {
	fd := int(lb.lbr.fd.Fd())
	responseMsg := Message{
		Id:      "go::1",
		Data:    data,
		MsgType: msgType,
	}
	jsonData, _ := json.Marshal(responseMsg)
	return syscall.Sendmsg(fd, append(jsonData, '\n'), nil, nil, 0)
}

func (lb *Bridge) listen() {
	for {
		lb.tryGetMessage()
		// execCount, _ := lb.tryGetMessage()
		// if execCount <= 0 {
		// 	time.Sleep(time.Microsecond * 100)
		// }
	}
}

func (lb *Bridge) tryGetMessage() (int, error) {
	data, err := lb.lbr.reader.ReadBytes('\n')
	if err != nil {
		return 0, err
	}
	msg := new(Message)
	json.Unmarshal(data, msg)
	event := msg.MsgType
	if event == "ready" {
		event = "establish"
	} else if event == "close" {
		lb.closeChan <- 1
	} else if !lb.established {
		return 0, nil
	}
	eventListener, exists := lb.eventMap[event]
	if !exists {
		return 0, nil
	}
	go eventListener(Context{
		Id:    msg.Id,
		Data:  msg.Data,
		lbRef: lb,
	})
	if event == "establish" {
		lb.established = true
	}
	return 1, nil
}

func (ctx *Context) Response(data string) error {
	fd := int(ctx.lbRef.lbr.fd.Fd())
	responseMsg := Message{
		Id:      ctx.Id,
		Data:    data,
		MsgType: "response",
	}
	jsonData, _ := json.Marshal(responseMsg)
	return syscall.Sendmsg(fd, append(jsonData, '\n'), nil, nil, 0)
}
