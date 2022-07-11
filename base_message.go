package godelayqueueclient

import "fmt"

var MessageAuthCode = "0_ONMARS_1"

type Command uint

const (
	Test Command = iota + 1
	Push
	Update
	Delete
)

type MessageBuilder interface {
	Build() string
}

type base struct {
	AuthCode string  `json:"auth_code"`
	CmdName  Command `json:"cmd_name"`
}

func (m *base) toStringArray() []string {
	var lines []string
	lines = append(lines, m.AuthCode)
	lines = append(lines, fmt.Sprintf("%d", m.CmdName))
	return lines
}
