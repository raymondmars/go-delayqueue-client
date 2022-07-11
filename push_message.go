package godelayqueueclient

import (
	"fmt"
	"strings"
)

var messageAuthCode = "0_ONMARS_1"

type pushMessage struct {
	AuthCode     string `json:"auth_code"`
	CmdName      string `json:"cmd_name"`
	DelaySeconds int64  `json:"delay_seconds"`
	TaskMode     string `json:"task_mode"`
	HttpUrl      string `json:"http_url"`
	Contents     string `json:"contents"`
}

func (m *pushMessage) BuildHttpModeMessage(delaySeconds int64, callbackUrl, contents string) string {
	m.AuthCode = messageAuthCode
	m.CmdName = "2"
	m.DelaySeconds = delaySeconds
	m.TaskMode = "1"
	m.HttpUrl = callbackUrl
	m.Contents = contents

	return m.toString()
}

func (m *pushMessage) BuildPingMessage() string {
	m.AuthCode = messageAuthCode
	m.CmdName = "1"
	return m.toString()
}

func (m *pushMessage) toString() string {
	var lines []string
	lines = append(lines, m.AuthCode)
	lines = append(lines, m.CmdName)
	if m.DelaySeconds > 0 {
		lines = append(lines, fmt.Sprintf("%d", m.DelaySeconds))
	}
	if m.TaskMode != "" {
		lines = append(lines, m.TaskMode)
	}
	if m.HttpUrl != "" {
		lines = append(lines, m.HttpUrl)
	}
	lines = append(lines, m.Contents)

	return strings.Join(lines, "\n")
}
