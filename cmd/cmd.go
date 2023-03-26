package cmd

import "strings"

type Type int

const (
	USER Type = iota
	QUIT
	PORT
	STRU
	RETR
	NOOP
)

func (t Type) String() string {
	var str string
	switch t {
	case USER:
		str = "USER"
	case QUIT:
		str = "QUIT"
	case PORT:
		str = "PORT"
	case RETR:
		str = "RETR"
	case NOOP:
		str = "NOOP"
	default:
		panic("There is a const having no String method.")
	}
	return str
}

func Parse(s string) (string, []string) {
	args := strings.Split(strings.Trim(s, " "), " ")
	cmd := strings.ToUpper(args[0])
	args = args[1:]
	return cmd, args
}
