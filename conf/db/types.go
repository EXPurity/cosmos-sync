package db

import (
	"cosmos-sync/logger"
	"cosmos-sync/util/constant"
	"os"
)

var (
	Addrs    = "192.168.150.7:30000"
	User     = "cosmos"
	Passwd   = "cosmos"
	Database = "sync-cosmos"
)

// get value of env var
func init() {
	addrs, found := os.LookupEnv(constant.EnvNameDbAddr)
	if found {
		Addrs = addrs
	}
	logger.Info("Env Value", logger.String(constant.EnvNameDbAddr, Addrs))

	user, found := os.LookupEnv(constant.EnvNameDbUser)
	if found {
		User = user
	}
	logger.Info("Env Value", logger.String(constant.EnvNameDbUser, User))

	passwd, found := os.LookupEnv(constant.EnvNameDbPassWd)
	if found {
		Passwd = passwd
	}
	logger.Info("Env Value", logger.String(constant.EnvNameDbPassWd, Passwd))

	database, found := os.LookupEnv(constant.EnvNameDbDataBase)
	if found {
		Database = database
	}
	logger.Info("Env Value", logger.String(constant.EnvNameDbDataBase, Database))
}
