package godelayqueueclient

import (
	"fmt"
	"strings"
)

type pingMessage struct {
	base
}

func NewPingMessage() MessageBuilder {
	p := &pingMessage{}
	p.AuthCode = MessageAuthCode
	p.CmdName = Test
	return p
}
func (p *pingMessage) Build() string {
	fmt.Println("-========================")
	fmt.Println(p)
	contents := fmt.Sprintf("%s\n\n", strings.Join(p.toStringArray(), "\n"))
	return contents
}
