package main

import (
	"log"
	"os"
	"wannabe/config"
	"wannabe/handlers"
	"wannabe/providers"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	// TODO read config path from env variable
	config := config.Load("config.json")

	storageProvider, err := providers.StorageProviderFactory(config.StorageProvider)
	if err != nil {
		log.Fatalf("fatal error when starting app: %v", err)
	}

	app := fiber.New()

	// TODO implement logger using factory pattern identical to StorageProviderFactory
	// Initialize logger
	if config.Logger.Enabled {
		file, err := os.OpenFile(config.Logger.Filepath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("fatal error when starting app: %v", err)
		}

		defer file.Close()

		app.Use(logger.New(logger.Config{
			Output: file,
			Format: config.Logger.Format,
		}))
	}

	// Probes endpoints
	app.Get("/wannabe/liveness", handlers.Liveness)
	app.Get("/wannabe/readiness", handlers.Readiness)

	// Api endpoints
	// NOTE
	app.Get("/wannabe/api/records/:hash", handlers.GetRecords(storageProvider))
	app.Get("/wannabe/api/records", handlers.GetRecords(storageProvider))
	app.Post("/wannabe/api/records", handlers.PostRecords(config, storageProvider))
	app.Delete("/wannabe/api/records/:hash", handlers.DeleteRecords(storageProvider))
	app.Delete("/wannabe/api/records", handlers.DeleteRecords(storageProvider))

	// NOTE get all hashes/curls
	// app.Get("/wannabe/api/hash/curl", handlers.Api)

	// NOTE stats
	// app.Get("/wannabe/api/stats", handlers.Api)
	// NOTE logs
	// app.Get("/wannabe/api/logs", handlers.Api)
	// app.Delete("/wannabe/api/logs", handlers.Api)

	// NOTE run all requests and validate responses
	// app.Get("/wannabe/api/records/test", handlers.Api)

	// Tools endpoints
	// NOTE validate config
	// app.Get("/wannabe/tools/validate-config", handlers.Tools)
	// NOTE regenerate records using new config
	// app.Get("/wannabe/tools/regenerate", handlers.Tools)

	// Wannabe endpoints
	app.Get("/*", handlers.Wannabe(config, storageProvider))
	app.Post("/*", handlers.Wannabe(config, storageProvider))

	// TODO read host and port from env variable
	app.Listen("localhost:1234")
}
