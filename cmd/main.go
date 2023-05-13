package main

import (
	"flag"
	"github.com/Natali-Skv/ProtectMyPassBot/config"
	passmanProto "github.com/Natali-Skv/ProtectMyPassBot/internal/passman/proto"
	"github.com/Natali-Skv/ProtectMyPassBot/internal/telegram_bot/delivery"
	tgBotRepo "github.com/Natali-Skv/ProtectMyPassBot/internal/telegram_bot/repository"
	tgBotUcase "github.com/Natali-Skv/ProtectMyPassBot/internal/telegram_bot/usecase"
	tarantoolTool "github.com/Natali-Skv/ProtectMyPassBot/internal/tools/tarantool"
	"github.com/tarantool/go-tarantool"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

func main() {
	loggerConfig := zap.NewDevelopmentConfig()
	loggerConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	logger, err := loggerConfig.Build()
	if err != nil {
		log.Fatal("zap logger build error")
	}
	defer func(logger *zap.Logger) {
		err = logger.Sync()
		if err != nil {
			logger.Error("logger sync error", zap.Error(err))
		}
	}(logger)

	configPath := *flag.String("config", config.TgBotDefaultConfigPath, "path to config file")
	flag.Parse()

	tgbotConfig := config.BotConfig{}
	err = config.ReadConfig(configPath, &tgbotConfig)
	if err != nil {
		logger.Fatal("reading config error", zap.Error(err))
	}

	logger.Debug("config", zap.Any("config", tgbotConfig))

	passmanGrpcConn, err := grpc.Dial(
		tgbotConfig.PassmanAddr,
		// grpc.WithInsecure(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		logger.Fatal("error connecting to grpc passman microservice", zap.Error(err))
	}
	defer func(passmanGrpcConn *grpc.ClientConn) {
		err := passmanGrpcConn.Close()
		if err != nil {
			logger.Error("error closing grpc connection to passman microservice", zap.Error(err))
		}
	}(passmanGrpcConn)

	passmanCli := passmanProto.NewPassmanServiceClient(passmanGrpcConn)

	conn, err := tarantoolTool.NewTarantoolConn(tgbotConfig.Tarantool)
	if err != nil {
		logger.Fatal("error opening connection to tarantool", zap.Error(err))
	}
	defer func(conn *tarantool.Connection) {
		err := conn.Close()
		if err != nil {
			logger.Error("error closing connection to tarantool", zap.Error(err))
		}
	}(conn)

	tgbRepo := tgBotRepo.NewTgBotRepo(logger, conn)

	tgbUcase := tgBotUcase.NewTgBotUsecase(logger, tgbRepo, passmanCli)

	tgBot := delivery.NewTelegramBot(tgbotConfig.Bot, logger, tgbUcase)
	err = tgBot.Run()
	if err != nil {
		logger.Fatal("running telegram bot error", zap.Error(err))
	}
}

//
//func handleSetCommand(message *tgbotapi.Message, client *tarantool.Connection) string {
//	args := message.CommandArguments()
//	if len(args) == 0 {
//		return "Usage: /set <service> <login> <password>"
//	}
//
//	parts := strings.Split(args, " ")
//	if len(parts) != 3 {
//		return "Usage: /set <service> <login> <password>"
//	}
//
//	service := parts[0]
//	login := parts[1]
//	password := parts[2]
//
//	_, err := client.Insert("passwords", []interface{}{message.From.ID, service, login, password})
//	if err != nil {
//		return fmt.Sprintf("Error setting password: %v", err)
//	}
//
//	return "password saved"
//}
//
//func handleGetCommand(message *tgbotapi.Message, client *tarantool.Connection) string {
//	args := message.CommandArguments()
//	if len(args) == 0 {
//		return "Usage: /get <service>"
//	}
//
//	service := args
//
//	res, err := client.Eval("return box.space.passwords.index.secondary:select({?, ?})", []interface{}{message.From.ID, service})
//	if err != nil {
//		return fmt.Sprintf("Error getting password: %v", err)
//	}
//
//	if len(res.Data) == 0 {
//		return "password not found"
//	}
//
//	tuple, ok := res.Data[0].([]interface{})
//	if !ok {
//		return "Error getting password"
//	}
//
//	login, ok := tuple[2].(string)
//	if !ok {
//		return "Error getting password"
//	}
//
//	password, ok := tuple[3].(string)
//	if !ok {
//		return "Error getting password"
//	}
//
//	return fmt.Sprintf("login: %s\nPassword: %s", login, password)
//}
//
//func handleDelCommand(message *tgbotapi.Message, client *tarantool.Connection) string {
//	args := message.CommandArguments()
//	if len(args) == 0 {
//		return "Usage: /del <service>"
//	}
//
//	service := args
//
//	res, err := client.Eval("return box.space.passwords.index.secondary:select({?,?})", []interface{}{message.From.ID, service})
//	if err != nil {
//		return fmt.Sprintf("Error deleting password: %v", err)
//	}
//
//	if len(res.Data) == 0 {
//		return "password not found"
//	}
//
//	tuple := res.Data[0].([]interface{})
//	_, err = client.Delete("passwords", "primary", tuple[0])
//	if err != nil {
//		return fmt.Sprintf("Error deleting password: %v", err)
//	}
//
//	return "password deleted"
//}
//
////func handleDelCommand(message *tgbotapi.Message, client *tarantool.Connection) string {
////	args := message.CommandArguments()
////	if len(args) == 0 {
////		return "Usage: /del <service>"
////	}
////
////	service := args
////
////	// Select the tuple to delete
////	resp, err := client.Select("passwords", "secondary", 0, 1, tarantool.IterEq, []interface{}{message.From.ID, service})
////	if err != nil {
////		return fmt.Sprintf("Error deleting password: %v", err)
////	}
////
////	if len(resp.Tuples()) == 0 {
////		return "password not found"
////	}
////
////	tuple := resp.Tuples()[0]
////
////	// Delete the tuple
////	_, err = client.Delete("passwords", tuple[0])
////	if err != nil {
////		return fmt.Sprintf("Error deleting password: %v", err)
////	}
////
////	return "password deleted"
////}
//
//func deleteMessageAfterDelay(bot *tgbotapi.BotAPI, message *tgbotapi.Message, delay time.Duration) {
//	time.Sleep(delay)
//
//	deleteMsg := tgbotapi.NewDeleteMessage(message.Chat.ID, message.MessageID)
//	_, err := bot.DeleteMessage(deleteMsg)
//	if err != nil {
//		log.Printf("Error deleting message: %v", err)
//	}
//}
