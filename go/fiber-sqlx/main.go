package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/ashwins93/fiber-sql/db"
	"github.com/ashwins93/fiber-sql/routes"
	"github.com/ashwins93/fiber-sql/utils"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	dbUrl := os.Getenv("GO_DB_URL")
	if dbUrl == "" {
		log.Fatal("GO_DB_URL is not set")
	}

	conn, err := sqlx.Connect("sqlite3", fmt.Sprintf("file:%s?_fk=1", dbUrl))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	queries := db.NewDb(conn)
	app := fiber.New(fiber.Config{
		JSONEncoder:  json.Marshal,
		JSONDecoder:  json.Unmarshal,
		ErrorHandler: errorHandler,
	})
	app.Use(recover.New())
	app.Use(logger.New())
	idGen := utils.NewNanoIDGenerator(21)

	server := routes.NewService(queries, app, idGen)
	server.SetupV1Routes()

	app.Hooks().OnShutdown(func() error {
		return conn.Close()
	})

	if err = app.Listen(":3000"); err != nil {
		log.Fatal(err)
	}

}

func errorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	message := "Something went wrong"

	var e *fiber.Error
	if errors.As(err, &e) {
		code = e.Code
		message = e.Message
	}

	return c.Status(code).JSON(&fiber.Map{
		"message": message,
	})
}
