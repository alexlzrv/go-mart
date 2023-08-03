package main

import (
	"log"

	apigomart "github.com/alexlzrv/go-mart/internal/api-go-mart"
	"github.com/alexlzrv/go-mart/internal/config"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}
	apigomart.Run(cfg)
}
