package telegram_bot

type Command struct {
	Name          string
	Usage         string
	ArgumentCount int
	RespFmtString string
}

var (
	HelpCommand  = Command{Name: "help", Usage: "Usage: `/help`", ArgumentCount: 0, RespFmtString: "*Set login and password to service:*\n`/set <service> <login> <password>`\n\n*Get credentials by service name:*\n`/get <service>`\n\n*Delete login and password by service name:*\n`/del <service>`\n\n*Help*\n/help"}
	StartCommand = Command{Name: "start", Usage: "Usage: `/start`", ArgumentCount: 0, RespFmtString: "*Set login and password to service:*\n`/set <service> <login> <password>`\n\n*Get credentials by service name:*\n`/get <service>`\n\n*Delete login and password by service name:*\n`/del <service>`\n\n*Help*\n/help"}
	SetCommand   = Command{Name: "set", Usage: "Usage: `/set <service> <login> <password>`", ArgumentCount: 3}
	DelCommand   = Command{Name: "del", Usage: "Usage: `/del <service>`", ArgumentCount: 1}
	GetCommand   = Command{Name: "get", Usage: "Usage: `/get <service>`", ArgumentCount: 1, RespFmtString: "Service: `%s` \nLogin: `%s`\nPassword: `%s`"}

	GetCommandServiceArgumentNumber = 0
)

const (
	UnknownErrorResp = "Sorry, unknown error occurred, try again later."
	NoSuchCredsMsg   = "No such credentials, try another service name"
	NoSuchUserMsg    = "Sorry, you do not have any saved credentials"
)
