package enum

type ClientCommand string

const (
	ClientCommandAttack ClientCommand = "attack"
	ClientCommandLogin  ClientCommand = "login"
	ClientCommandLogout ClientCommand = "logout"
)
