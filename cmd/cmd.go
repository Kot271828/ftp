package cmd

import "strings"

type Type int

const (
	USER Type = iota
	PWD
	LIST
	QUIT
	PORT
	TYPE
	MODE
	STRU
	RETR
	NOOP
	UNKNOWN
)

func (t Type) String() string {
	var str string
	switch t {
	case USER:
		str = "USER"
	case PWD:
		str = "PWD"
	case LIST:
		str = "LIST"
	case QUIT:
		str = "QUIT"
	case TYPE:
		str = "TYPE"
	case MODE:
		str = "MODE"
	case STRU:
		str = "STRU"
	case PORT:
		str = "PORT"
	case RETR:
		str = "RETR"
	case NOOP:
		str = "NOOP"
	case UNKNOWN:
		str = "UNKNOWN"
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
	case "LIST":
		cmd = LIST
	case "RETR":
		cmd = RETR
	case "TYPE":
		cmd = TYPE
	case "MODE":
		cmd = MODE
	case "STRU":
		cmd = STRU
	case "PORT":
		cmd = PORT
	default:
		cmd = UNKNOWN
	}
	args = args[1:]
	return cmd, args
}
