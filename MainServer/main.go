package main

import (
	"Crypto_Bot/Scrapper/server"
	"os"
)

func main() {
	server_url, exist := os.LookupEnv("SERVER_HOST")
	if !exist {
		return
	}
	db, exist := os.LookupEnv("DB_URL")
	if !exist {
		return
	}
	server, err := server.BuildServer(server_url, db)
	if err != nil {
		return
	}
	defer server.Stop()
	server.Start()
}
