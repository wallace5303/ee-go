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
	br          *BridgeReader
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
	bRef    *Bridge
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
func (br *BridgeReader) Read(p []byte) (n int, err error) {
	// 后去传递的数据的有效负载长度
	length := syscall.CmsgSpace(4)
	mbuf := make([]byte, length)
	// 将收到的数据从内核空间拷贝至用户空间
	// oob: Out of band data https://beej.us/298C/oob_overview.html
	n, _, _, _, err = syscall.Recvmsg(int(br.fd.Fd()), p, mbuf, 0)
	if err != nil {
		return 0, err
	}
	return n, nil
}

func Establish() *Bridge {
	reader := establishChannel()
	b := &Bridge{
		br:       reader,
		eventMap: make(map[string]EventCallback),
	}
	return b
}

func establishChannel() *BridgeReader {
	nodeChannelFD := os.Getenv(NODE_CHANNEL_FD)
	nodeChannelFDInt, _ := strconv.Atoi(nodeChannelFD)
	fd := os.NewFile(uintptr(int(nodeChannelFDInt)), "lbipc"+nodeChannelFD)
	reader := &BridgeReader{fd: fd}
	reader.reader = bufio.NewReader(reader)
	return reader
}

func (b *Bridge) On(event string, callback EventCallback) {
	b.eventMap[event] = callback
}

func (b *Bridge) Listen() {
	b.closeChan = make(chan int)
	go func() {
		b.listen()
	}()
	b.sendByType("ready", "establish")
	<-b.closeChan
}

func (b *Bridge) sendByType(data string, msgType string) error {
	fd := int(b.br.fd.Fd())
	responseMsg := Message{
		Id:      "go::1",
		Data:    data,
		MsgType: msgType,
	}
	jsonData, _ := json.Marshal(responseMsg)
	return syscall.Sendmsg(fd, append(jsonData, '\n'), nil, nil, 0)
}

func (b *Bridge) listen() {
	for {
		b.tryGetMessage()
		// execCount, _ := b.tryGetMessage()
		// if execCount <= 0 {
		// 	time.Sleep(time.Microsecond * 100)
		// }
	}
}

func (b *Bridge) tryGetMessage() (int, error) {
	data, err := b.br.reader.ReadBytes('\n')
	if err != nil {
		return 0, err
	}
	msg := new(Message)
	json.Unmarshal(data, msg)
	event := msg.MsgType
	if event == "ready" {
		event = "establish"
	} else if event == "close" {
		b.closeChan <- 1
	} else if !b.established {
		return 0, nil
	}
	eventListener, exists := b.eventMap[event]
	if !exists {
		return 0, nil
	}
	go eventListener(Context{
		Id:   msg.Id,
		Data: msg.Data,
		bRef: b,
	})
	if event == "establish" {
		b.established = true
	}
	return 1, nil
}

func (ctx *Context) Response(data string) error {
	fd := int(ctx.bRef.br.fd.Fd())
	responseMsg := Message{
		Id:      ctx.Id,
		Data:    data,
		MsgType: "response",
	}
	jsonData, _ := json.Marshal(responseMsg)
	return syscall.Sendmsg(fd, append(jsonData, '\n'), nil, nil, 0)
}
