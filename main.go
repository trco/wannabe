package main

import (
	"log"
	"wannabe/config"
	"wannabe/handlers"
	"wannabe/providers"

	"github.com/gofiber/fiber/v2"
)

func main() {
	// TODO read config path from env variable
	config, err := config.Load("config.json")
	if err != nil {
		log.Fatalf("fatal error starting app: %v", err)
	}

	storageProvider, err := providers.StorageProviderFactory(config.StorageProvider)
	if err != nil {
		log.Fatalf("fatal error starting app: %v", err)
	}

	app := fiber.New()

	// Probes endpoints
	app.Get("/wannabe/liveness", handlers.Liveness)
	app.Get("/wannabe/readiness", handlers.Readiness)

	// Api endpoints
	app.Get("/wannabe/api/records/:hash", handlers.GetRecords(storageProvider))
	app.Get("/wannabe/api/records", handlers.GetRecords(storageProvider))
	app.Post("/wannabe/api/records", handlers.PostRecords(config, storageProvider))
	app.Delete("/wannabe/api/records/:hash", handlers.DeleteRecords(storageProvider))
	app.Delete("/wannabe/api/records", handlers.DeleteRecords(storageProvider))
	app.Get("/wannabe/api/regenerate", handlers.Regenerate(app, config, storageProvider))

	// Wannabe endpoints
	app.Get("/*", handlers.Wannabe(config, storageProvider))
	app.Post("/*", handlers.Wannabe(config, storageProvider))

	// TODO read host and port from env variable
	app.Listen("localhost:1234")
}
