package uf_aim

import (
	"strings"
)

func (aimInfo *AimInfo) IsShortcutQuit() bool {
	if aimInfo.IsShortcut() {
		cmd := aimInfo.GetMsgIn()
		cmd = strings.ToLower(cmd)
		cmd = strings.ReplaceAll(cmd, " ", "")
		cmd = strings.ReplaceAll(cmd, "\n", "")

		switch cmd {
		case "cancel":
			fallthrough
		case "stop":
			fallthrough
		case "quit":
			return true
		}
	}
	return false
}

func IsShortcut(in string) bool {
	if len(in) > longestShortcut {
		return false
	}

	inLower := strings.ToLower(in)
	inLower = strings.ReplaceAll(inLower, " ", "")
	inLower = strings.ReplaceAll(inLower, "\n", "")
	for _, shortcut := range treeShortcuts {
		if inLower == shortcut {
			return true
		}
	}
	return false
}

var longestShortcut = 4

var treeShortcuts = []string{
	"cancel",
	"hi",
	"hey",
	"hell",
	"hello",
	"howdy",
	"n",
	"naw",
	"no",
	"no!",
	"quit",
	"stop",
	"y",
	"ye",
	"yes",
	"yes!",
	"q",
	"1",
	"2",
	"3",
	"4",
	"5",
	"6",
	"7",
	"8",
	"9",
	"0",
}
