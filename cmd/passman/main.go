package main

import (
	"flag"
	config2 "github.com/Natali-Skv/ProtectMyPassBot/internal/config"
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

	configPath := flag.String("config", config2.PassmanDefaultConfigPath, "path to config file")
	flag.Parse()

	passmanConfig := config2.PassmanConfig{}
	err = config2.ReadConfig(*configPath, &passmanConfig)
	if err != nil {
		logger.Fatal("reading config error", zap.Error(err))
	}

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
	usecase := passmahUsecase.NewPassmanUsecase(logger, repo)
	handler := passmanHandler.NewPassmanHandler(logger, usecase)

	passmanProto.RegisterPassmanServiceServer(server, handler)

	logger.Info("passman microservice is starting")

	err = server.Serve(lis)

	if err != nil {
		log.Fatalln("cant serve passman microservice", err)
	}
}
