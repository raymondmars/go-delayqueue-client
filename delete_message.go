package godelayqueueclient

import (
	"fmt"
	"strings"
)

type deleteMessage struct {
	base
	TaskId string `json:"task_id"`
}

func NewDeleteMessage(taskId string) MessageBuilder {
	dm := &deleteMessage{TaskId: taskId}
	dm.AuthCode = MessageAuthCode
	dm.CmdName = Delete
	return dm
}

func (d *deleteMessage) Build() string {
	lines := d.toStringArray()
	lines = append(lines, d.TaskId)
	contents := fmt.Sprintf("%s\n\n", strings.Join(lines, "\n"))
	return contents
}
