package main

import (
	"flag"

	"github.com/localhots/cmdui/backend/api"
	"github.com/localhots/cmdui/backend/commands"
	"github.com/localhots/cmdui/backend/config"
	"github.com/localhots/cmdui/backend/db"
	"github.com/localhots/cmdui/backend/log"
	"github.com/localhots/cmdui/backend/runner"
)

func main() {
	confPath := flag.String("config", "config.toml", "Path to config file in TOML format")
	flag.Parse()

	_, err := config.LoadFile(*confPath)
	if err != nil {
		log.WithFields(log.F{"error": err}).Fatal("Failed to load config file")
	}
	if err := db.Connect(); err != nil {
		log.WithFields(log.F{"error": err}).Fatal("Failed to establish database connection")
	}
	if err := runner.PrepareLogsDir(); err != nil {
		log.WithFields(log.F{"error": err}).Fatal("Failed to create logs directory")
	}
	list, err := runner.CommandsList()
	if err != nil {
		log.WithFields(log.F{"error": err}).Fatal("Failed to import commands")
	}
	commands.Import(list)
	defer runner.Shutdown()
	if err := api.Start(); err != nil {
		log.WithFields(log.F{"error": err}).Fatal("Failed to start the server")
	}
}
