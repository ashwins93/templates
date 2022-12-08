package routes

import (
	"github.com/ashwins93/fiber-badger/db"
	"github.com/ashwins93/fiber-badger/utils"
	"github.com/gofiber/fiber/v2"
)

func (s *Service) setupUserRoutes(router fiber.Router) {
	router.Post("", s.createUserHandler)
	router.Get("", s.findUsersHandler)
}

func (s *Service) createUserHandler(c *fiber.Ctx) error {
	userParams := db.CreateUserParams{}

	if err := c.BodyParser(&userParams); err != nil {
		return err
	}

	errors := utils.ValidateStruct(userParams)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errors)
	}

	user, err := s.queries.CreateNewUser(&userParams)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(user)
}

func (s *Service) findUsersHandler(c *fiber.Ctx) error {
	users, err := s.queries.GetAllUsers()
	if err != nil {
		return err
	}

	return c.JSON(users)
}
