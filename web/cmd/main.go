package main

import (
	"log"

	"github.com/H1K0/SkazaNull/conf"
	"github.com/H1K0/SkazaNull/db"
	"github.com/H1K0/SkazaNull/server"
)

func main() {
	config, err := conf.GetConf()
	if err != nil {
		log.Fatalf("error while loading configuration file: %s\n", err)
	}
	db.InitDB(config.DBConfig)
	server.Serve(config.Interface, []byte(config.EncryptionKey))
}
