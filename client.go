package godelayqueueclient

import (
	"fmt"
	"net"
	"time"

	log "github.com/sirupsen/logrus"
)

type Client interface {
	Push(futureDate time.Time, callbackUrl, contents string) error
	Ping() (error, string)
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

func (c *client) Push(futureDate time.Time, callbackUrl string, contents string) error {
	now := time.Now().UTC()
	diff := now.Sub(futureDate.UTC())
	remoteConn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%s", delayQueueHost, delayQueuePort), time.Second*5)
	if err != nil {
		return err
	}
	message := &pushMessage{}
	toMsg := message.BuildHttpModeMessage(int64(diff.Seconds()), callbackUrl, contents)
	log.Info(toMsg)
	defer remoteConn.Close()
	remoteConn.Write([]byte(toMsg))
	return nil
}

func (c *client) Ping() (error, string) {
	return nil, ""
}
