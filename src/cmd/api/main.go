package main

import "github.com/saifutdinov/go-invoices-api/api"

const (
	configPath = "config.toml"
)

func main() {
	api.RunServer(configPath)
}
