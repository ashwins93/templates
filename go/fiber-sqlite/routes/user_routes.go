package routes

import (
	"github.com/ashwins93/fiber-sql/db"
	"github.com/ashwins93/fiber-sql/utils"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

func (s *Service) setupUserRoutes(router fiber.Router) {
	router.Post("", s.createUserHandler)
	router.Get("", s.findUsersHandler)
	router.Get("/:id", s.findUserByIDHandler)
	router.Patch("/:id", s.updateUserHandler)
	router.Delete("/:id", s.deleteUserHandler)
}

func (s *Service) createUserHandler(c *fiber.Ctx) error {
	userParams := db.CreateUserParams{}

	if err := c.BodyParser(&userParams); err != nil {
		return err
	}

	userParams.ID = s.idGen.Generate()
	errors := utils.ValidateStruct(userParams)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errors)
	}

	existingUser, err := s.queries.GetUserByEmail(c.Context(), userParams.Email)
	if err != nil {
		return err
	}

	if existingUser.Email != "" {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"message": "Email already exists",
		})
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(userParams.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	userParams.Password = string(hash)

	user, err := s.queries.CreateUser(c.Context(), userParams)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(user)
}

func (s *Service) findUsersHandler(c *fiber.Ctx) error {
	users, err := s.queries.GetUsers(c.Context())
	if err != nil {
		return err
	}

	return c.JSON(users)
}

func (s *Service) findUserByIDHandler(c *fiber.Ctx) error {
	id := c.Params("id")
	user, err := s.queries.GetUserByID(c.Context(), id)
	if err != nil {
		return err
	}

	if user.Email == "" {
		return c.Status(fiber.StatusNotFound).JSON(&fiber.Map{
			"message": "User not found",
		})
	}

	return c.JSON(user)
}

func (s *Service) updateUserHandler(c *fiber.Ctx) error {
	id := c.Params("id")
	existingUser, err := s.queries.GetUserByID(c.Context(), id)
	if err != nil {
		return err
	}

	if existingUser.Email == "" {
		return c.Status(fiber.StatusNotFound).JSON(&fiber.Map{
			"message": "User not found",
		})
	}

	userParams := db.UpdateUserParams{}

	if err := c.BodyParser(&userParams); err != nil {
		return err
	}

	errors := utils.ValidateStruct(userParams)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errors)
	}

	user, err := s.queries.UpdateUser(c.Context(), userParams, id)
	if err != nil {
		return err
	}

	return c.JSON(user)
}

func (s *Service) deleteUserHandler(c *fiber.Ctx) error {
	id := c.Params("id")
	existingUser, err := s.queries.GetUserByID(c.Context(), id)
	if err != nil {
		return err
	}

	if existingUser.Email == "" {
		return c.Status(fiber.StatusNotFound).JSON(&fiber.Map{
			"message": "User not found",
		})
	}

	err = s.queries.DeleteUser(c.Context(), id)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusNoContent).Send(nil)
}
