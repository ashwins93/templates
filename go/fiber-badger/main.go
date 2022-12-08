package main

import (
	"errors"
	"log"
	"time"

	"github.com/ashwins93/fiber-badger/db"
	"github.com/ashwins93/fiber-badger/routes"
	"github.com/dgraph-io/badger/v3"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	badger, err := badger.Open(badger.DefaultOptions("./badger"))
	if err != nil {
		log.Fatal(err)
	}
	defer badger.Close()
	s := db.NewDb(badger)

	app := fiber.New(fiber.Config{
		JSONEncoder:  json.Marshal,
		JSONDecoder:  json.Unmarshal,
		ErrorHandler: errorHandler,
	})
	app.Use(recover.New())
	app.Use(logger.New())

	server := routes.NewService(s, app)

	server.SetupV1Routes()

	if err = app.Listen(":3000"); err != nil {
		log.Fatal(err)
	}

	app.Hooks().OnShutdown(func() error {
		return badger.Close()
	})

	go gcRunner(badger)
}

func gcRunner(db *badger.DB) {
	ticker := time.NewTicker(time.Minute * 5)
	defer ticker.Stop()

	for range ticker.C {
		var err error
		for err == nil {
			err = db.RunValueLogGC(0.7)
		}
	}
}

func errorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError

	var e *fiber.Error
	if errors.As(err, &e) {
		code = e.Code
	}

	return c.Status(code).JSON(&fiber.Map{
		"message": e.Message,
	})
}
