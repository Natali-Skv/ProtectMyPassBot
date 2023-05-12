package tarantool

import (
	"errors"
	"github.com/Natali-Skv/ProtectMyPassBot/config"
	"github.com/tarantool/go-tarantool"
)

var (
	ErrorCreatingConnection = errors.New("error creating connection to database")
)

func NewTarantoolConn(conf config.TarantoolConfig) (conn *tarantool.Connection, err error) {
	opts := tarantool.Opts{User: conf.User, Pass: conf.Password}
	conn, err = tarantool.Connect(conf.Address, opts)
	if err != nil {
		return nil, errors.Join(ErrorCreatingConnection, err)
	}
	return conn, nil
}
