package main

import (
	"flag"
	"github.com/Natali-Skv/ProtectMyPassBot/config"
	passmanProto "github.com/Natali-Skv/ProtectMyPassBot/internal/passman/proto"
	"github.com/Natali-Skv/ProtectMyPassBot/internal/telegram_bot/delivery"
	tgBotRepo "github.com/Natali-Skv/ProtectMyPassBot/internal/telegram_bot/repository"
	tgBotUcase "github.com/Natali-Skv/ProtectMyPassBot/internal/telegram_bot/usecase"
	"github.com/Natali-Skv/ProtectMyPassBot/internal/tools/delay_task_manager/delay_task_manager"
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

	passmanGrpcConn, err := grpc.Dial(
		tgbotConfig.PassmanAddr,
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

	delayTaskManager := delay_task_manager.NewDelayTaskManager(logger)
	tgBot := delivery.NewTelegramBot(tgbotConfig.Bot, logger, tgbUcase, delayTaskManager)
	err = tgBot.Run()
	if err != nil {
		logger.Fatal("running telegram bot error", zap.Error(err))
	}
}
