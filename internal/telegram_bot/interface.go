package telegram_bot

type Command struct {
	Name          string
	Usage         string
	ArgumentCount int
	RespFmtString string
}

var (
	SetCommand = Command{Name: "set", Usage: "Usage: `/set <service> <login> <password>`", ArgumentCount: 3}
	DelCommand = Command{Name: "del", Usage: "Usage: `/del <service>`", ArgumentCount: 1}
	GetCommand = Command{Name: "get", Usage: "Usage: `/get <service>`", ArgumentCount: 1, RespFmtString: "Service: `%s` \nLogin: `%s`\nPassword: `%s`"}

	GetCommandServiceArgumentNumber = 0
)

const (
	UnknownErrorResp = "Sorry, unknown error occurred, try again later."
	NoSuchCredsMsg   = "No such credentials, try another service name"
)
