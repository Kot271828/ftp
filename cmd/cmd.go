package cmd

import "strings"

type Type int

const (
	USER Type = iota
	PWD
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
	case PWD:
		str = "PWD"
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

func Parse(s string) (Type, []string) {
	args := strings.Split(strings.Trim(s, " "), " ")
	var cmd Type
	switch strings.ToUpper(args[0]) {
	case "QUIT":
		cmd = QUIT
	case "USER":
		cmd = USER
	case "PWD":
		cmd = PWD
	default:
		cmd = NOOP
	}
	args = args[1:]
	return cmd, args
}
