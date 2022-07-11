package godelayqueueclient

import (
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

type ResponseStatusCode uint
type ResponseErrCode uint

const (
	Ok ResponseStatusCode = iota + 1
	Fail
)

type NotifyMode uint

const (
	HTTP NotifyMode = iota + 1
	SubPub
)

const (
	NOT_READY            ResponseErrCode = 1000
	INVALID_MESSAGE      ResponseErrCode = 1010
	AUTH_FAILED          ResponseErrCode = 1012
	INVALID_DELAY_TIME   ResponseErrCode = 1014
	INVALID_COMMAND      ResponseErrCode = 1016
	INVALID_PUSH_MESSAGE ResponseErrCode = 1018
	UPDATE_FAILED        ResponseErrCode = 1020
	DELETE_FAILED        ResponseErrCode = 1022
)

type Response struct {
	Status    ResponseStatusCode
	ErrorCode ResponseErrCode
	Message   string
	TaskId    string
}

type Client interface {
	Push(delaySeconds int64, taskMode NotifyMode, taskTarget, contents string) (*Response, error)
	Delete(taskId string) (*Response, error)
	Ping() (*Response, error)
}

var delayQueueHost = getEvnWithDefaultVal("DELAY_QUEUE_HOST", "127.0.0.1")
var delayQueuePort = getEvnWithDefaultVal("DELAY_QUEUE_PORT", "3450")

type client struct {
	Host string
	Port string
}

func NewClient() Client {
	return &client{
		Host: delayQueueHost,
		Port: delayQueuePort,
	}
}

func NewClientWithHostAndPort(host, port string) Client {
	return &client{
		Host: host,
		Port: port,
	}
}

func (c *client) Push(delaySeconds int64, taskMode NotifyMode, taskTarget, contents string) (*Response, error) {
	builder := NewPushMessage(delaySeconds, strconv.FormatUint(uint64(taskMode), 10), taskTarget, contents)
	return c.send(builder)
}

func (c *client) Ping() (*Response, error) {
	return c.send(NewPingMessage())
}

func (c *client) Delete(taskId string) (*Response, error) {
	return c.send(NewDeleteMessage(taskId))
}

func (c *client) getConnection() (net.Conn, error) {
	return net.DialTimeout("tcp", fmt.Sprintf("%s:%s", c.Host, c.Port), time.Second*5)
}

func (c *client) send(builder MessageBuilder) (*Response, error) {
	remoteConn, err := c.getConnection()
	if err != nil {
		return &Response{}, err
	}
	toMsg := builder.Build()
	log.Info(toMsg)
	defer remoteConn.Close()
	_, err = remoteConn.Write([]byte(toMsg))
	if err != nil {
		return &Response{}, err
	}
	return c.buildResponseMessage(remoteConn)
}

func (c *client) buildResponseMessage(remoteConn net.Conn) (*Response, error) {
	reply := make([]byte, 1024)
	_, err := remoteConn.Read(reply)
	if err != nil {
		return &Response{}, err
	}
	serverResp := strings.Split(string(reply), "|")
	if len(serverResp) == 2 {
		return &Response{Status: Ok, TaskId: serverResp[1]}, nil
	} else if len(serverResp) == 3 {
		code, _ := strconv.Atoi(serverResp[1])
		return &Response{Status: Fail}, errors.New(fmt.Sprintf("Err Code: %d, %s", ResponseErrCode(code), serverResp[2]))
	} else {
		return &Response{}, errors.New(fmt.Sprintf("invalid server response: %s", serverResp))
	}
}
