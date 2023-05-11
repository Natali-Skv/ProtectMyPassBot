package main

import (
	"flag"
	"github.com/Natali-Skv/ProtectMyPassBot/config"
	passmanHandler "github.com/Natali-Skv/ProtectMyPassBot/internal/passman/delivery/grpc"
	passmanProto "github.com/Natali-Skv/ProtectMyPassBot/internal/passman/proto"
	passmahUsecase "github.com/Natali-Skv/ProtectMyPassBot/internal/passman/usecase"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"log"
	"net"
)

const (
	defaultConfigPath = "config/config.yaml"
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

	usecase := passmahUsecase.PassmanUsecase{}
	handler := passmanHandler.NewPassmanHandler(&usecase)

	passmanProto.RegisterPassmanServiceServer(server, handler)

	err = server.Serve(lis)

	if err != nil {
		log.Fatalln("cant serve auth-microservice", err)
	}
}
