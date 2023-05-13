package tarantool

import (
	"errors"
	"github.com/Natali-Skv/ProtectMyPassBot/config"
	"github.com/tarantool/go-tarantool"
	"time"
)

var (
	ErrorCreatingConnection = errors.New("error creating connection to database")
)

func NewTarantoolConn(conf config.TarantoolConfig) (conn *tarantool.Connection, err error) {
	opts := tarantool.Opts{
		User:          conf.User,
		Pass:          conf.Password,
		Timeout:       time.Duration(conf.Timeout) * time.Millisecond,
		Reconnect:     time.Duration(conf.Reconnect) * time.Second,
		MaxReconnects: conf.MaxReconnects,
	}
	conn, err = tarantool.Connect(conf.Address, opts)
	if err != nil {
		return nil, errors.Join(ErrorCreatingConnection, err)
	}
	return conn, nil
}
