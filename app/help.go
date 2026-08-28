package app

import "strings"

func Help() string {
	return strings.Join([]string{
		"toolnav sample",
		"toolnav list [category=<name>] [status=<name>]",
		"toolnav search <term>",
		"toolnav move <id> <position>",
		"toolnav backup [label]",
		"toolnav restore",
		"toolnav status <id> <active|beta|archived>",
		"toolnav dashboard",
		"toolnav audit",
		"toolnav export",
	}, "\n")
}
