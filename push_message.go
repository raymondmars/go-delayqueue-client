package godelayqueueclient

import (
	"fmt"
	"strings"
)

type pushMessage struct {
	base
	DelaySeconds int64  `json:"delay_seconds"`
	TaskMode     string `json:"task_mode"`
	// if taskmode is http, target is callbackURL
	// if taskmode is pubsub, target is queuename
	TaskTarget string `json:"task_target"`
	Contents   string `json:"contents"`
}

func NewPushMessage(delaySeconds int64, taskMode, taskTarget, contents string) MessageBuilder {
	p := &pushMessage{
		DelaySeconds: delaySeconds,
		TaskMode:     taskMode,
		TaskTarget:   taskTarget,
		Contents:     contents,
	}
	p.AuthCode = MessageAuthCode
	p.CmdName = Push
	return p
}

func (m *pushMessage) Build() string {
	lines := m.toStringArray()
	if m.DelaySeconds > 0 {
		lines = append(lines, fmt.Sprintf("%d", m.DelaySeconds))
	}
	if m.TaskMode != "" {
		lines = append(lines, m.TaskMode)
	}
	if m.TaskTarget != "" {
		lines = append(lines, m.TaskTarget)
	}
	if m.Contents != "" {
		lines = append(lines, m.Contents)
	}
	contents := fmt.Sprintf("%s\n\n", strings.Join(lines, "\n"))
	return contents
}
