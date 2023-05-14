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
	SetCommand   = Command{Name: "set", Usage: "Usage: `/set <service> <login> <password>`", ArgumentCount: 3, RespFmtString: "Your credentials for service `%s` saved."}
	DelCommand   = Command{Name: "del", Usage: "Usage: `/del <service>`", ArgumentCount: 1, RespFmtString: "Your credentials for service `%s` deleted."}
	GetCommand   = Command{Name: "get", Usage: "Usage: `/get <service>`", ArgumentCount: 1, RespFmtString: "Service: `%s` \nLogin: `%s`\nPassword: `%s`"}

	GetCommandServiceArgIdx = 0

	DelCommandServiceArgIdx = 0

	SetCommandServiceArgIdx  = 0
	SetCommandLoginArgIdx    = 1
	SetCommandPasswordArgIdx = 2
)

const (
	UnknownErrorResp = "Sorry, unknown error occurred, try again later."
	NoSuchCredsMsg   = "No such credentials, try another service name"
	NoSuchUserMsg    = "Sorry, you do not have any saved credentials"
)
