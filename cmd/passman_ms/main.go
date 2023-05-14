package main

import (
	"flag"
	"github.com/Natali-Skv/ProtectMyPassBot/config"
	passmanHandler "github.com/Natali-Skv/ProtectMyPassBot/internal/passman/delivery/grpc"
	passmanProto "github.com/Natali-Skv/ProtectMyPassBot/internal/passman/proto"
	passmahRepository "github.com/Natali-Skv/ProtectMyPassBot/internal/passman/repository"
	passmahUsecase "github.com/Natali-Skv/ProtectMyPassBot/internal/passman/usecase"
	tarantoolTool "github.com/Natali-Skv/ProtectMyPassBot/internal/tools/tarantool"
	"github.com/tarantool/go-tarantool"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"log"
	"net"
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
			log.Fatal(err)
		}
	}(logger)

	configPath := *flag.String("config", config.PassmanDefaultConfigPath, "path to config file")
	flag.Parse()

	passmanConfig := config.PassmanConfig{}
	err = config.ReadConfig(configPath, &passmanConfig)
	if err != nil {
		logger.Fatal("reading config error", zap.Error(err))
	}

	logger.Debug("passman microservice config", zap.Any("config", passmanConfig))

	lis, err := net.Listen("tcp", passmanConfig.BindAddr)
	if err != nil {
		logger.Fatal("can't listen port", zap.Error(err))
	}

	server := grpc.NewServer()

	conn, err := tarantoolTool.NewTarantoolConn(passmanConfig.Tarantool)
	if err != nil {
		logger.Fatal("error opening connection to tarantool", zap.Error(err))
	}
	defer func(conn *tarantool.Connection) {
		err := conn.Close()
		if err != nil {
			logger.Error("error closing connection to tarantool", zap.Error(err))
		}
	}(conn)

	repo := passmahRepository.NewPassmanRepo(logger, conn)

	// TODO delete

	userID, err := repo.Register()
	logger.Debug("register", zap.Error(err), zap.Int("useID", userID.Int()))
	//err = repo.AddCredentials(m.AddCredsReqR{UserID: m.UserID(1), Data: m.AddCredsData{
	//	Service:  "GG",
	//	Login:    "loginGG",
	//	Password: "passwordGG",
	//}})
	//logger.Debug("add creds", zap.Error(err))
	//
	//err = repo.DeleteCreds(m.DeleteCredsReqR{UserID: m.UserID(1), Service: "TG"})
	//logger.Debug("delete creds", zap.Error(err))

	usecase := passmahUsecase.NewPassmanUsecase(logger, repo)
	handler := passmanHandler.NewPassmanHandler(logger, usecase)

	passmanProto.RegisterPassmanServiceServer(server, handler)

	err = server.Serve(lis)

	if err != nil {
		log.Fatalln("cant serve passman microservice", err)
	}
}
