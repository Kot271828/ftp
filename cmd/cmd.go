package cmd

import "strings"

func Parse(s string) (string, []string) {
	args := strings.Split(strings.Trim(s, " "), " ")
	cmd := strings.ToUpper(args[0])
	args = args[1:]
	return cmd, args
}
