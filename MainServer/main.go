package main

import (
	"Crypto_Bot/MainServer/server"
	"os"
)

func main() {
	server_url, exist := os.LookupEnv("SERVER_HOST")
	if !exist {
		return
	}
	server := server.BuildServer(server_url, nil, nil, nil)
	defer server.Stop()
	server.Start()
}
