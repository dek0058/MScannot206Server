package framework

import (
	"fmt"
	"strings"
)

type ClientCommand interface {
	Commands() []string
	Execute(args []string) error
	Description() string
}

func MakeCommandDescription(cmds []string, argsHint string, description string) string {
	if len(cmds) == 0 {
		return ""
	}

	var sb strings.Builder
	// 첫 번째 명령어
	sb.WriteString(fmt.Sprintf("-%s", cmds[0]))
	if argsHint != "" {
		sb.WriteString(fmt.Sprintf(" %s", argsHint))
	}
	sb.WriteString(fmt.Sprintf(": %s", description))

	// 나머지 명령어 (별칭)
	for _, cmd := range cmds[1:] {
		sb.WriteString(fmt.Sprintf("\n-%s", cmd))
	}

	return sb.String()
}
