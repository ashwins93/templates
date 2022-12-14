package routes

import (
	"github.com/ashwins93/fiber-sql/db"
	"github.com/ashwins93/fiber-sql/utils"
	"github.com/gofiber/fiber/v2"
)

type Service struct {
	queries *db.Queries
	app     *fiber.App
	idGen   utils.IDGenerator
}

func NewService(queries *db.Queries, app *fiber.App, idGen utils.IDGenerator) *Service {
	return &Service{queries, app, idGen}
}

func (s *Service) SetupV1Routes() {
	v1Routes := s.app.Group("/api/v1")

	userRouter := v1Routes.Group("/users")

	s.setupUserRoutes(userRouter)
}
