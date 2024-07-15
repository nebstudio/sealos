package main

import (
	"log"

	"github.com/nebstudio/sealos/service/exceptionmonitor/api"
	"github.com/nebstudio/sealos/service/exceptionmonitor/dao"
	"github.com/nebstudio/sealos/service/exceptionmonitor/helper/client"
	"github.com/nebstudio/sealos/service/exceptionmonitor/helper/monitor"
)

func main() {
	if err := initialize(); err != nil {
		log.Printf("Initialization failed: %v", err)
	}
	go monitor.DatabaseExceptionMonitor()
	go monitor.DatabasePerformanceMonitor()
	go monitor.DatabaseBackupMonitor()
	select {}
}

func initialize() error {
	if err := api.GetENV(); err != nil {
		return err
	}
	if err := client.InitClient(); err != nil {
		return err
	}
	return dao.InitCockroachDB()
}
