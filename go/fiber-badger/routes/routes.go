package routes

import (
	"github.com/ashwins93/fiber-badger/db"
	"github.com/gofiber/fiber/v2"
)

type Service struct {
	queries *db.Queries
	app     *fiber.App
}

func NewService(queries *db.Queries, app *fiber.App) *Service {
	return &Service{queries, app}
}

func (s *Service) SetupV1Routes() {
	v1Routes := s.app.Group("/api/v1")

	userRouter := v1Routes.Group("/users")

	s.setupUserRoutes(userRouter)
}
