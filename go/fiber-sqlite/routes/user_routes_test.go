package routes

import (
	"bytes"
	"database/sql"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ashwins93/fiber-sql/db"
	"github.com/ashwins93/fiber-sql/utils"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/suite"
)

type UserRoutesTestSuite struct {
	suite.Suite
	q    *db.Queries
	conn *sql.DB
	app  *fiber.App
}

func (s *UserRoutesTestSuite) SetupSuite() {
	conn, err := sql.Open("sqlite3", "file:../db/test.db?_fk=1")
	if err != nil {
		panic(err)
	}

	s.q = db.NewDb(conn)
	s.conn = conn
	s.app = fiber.New(fiber.Config{
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
	})
	idGen := utils.NewNanoIDGenerator(21)

	service := NewService(s.q, s.app, idGen)
	service.SetupV1Routes()

	seedDataIntoDb(s.q)
}

func (s *UserRoutesTestSuite) BeforeTest(suite, testName string) {
	s.T().Log("BeforeTest: ", testName)

	switch testName {
	case "TestDeleteUser":
		s.conn.Exec("INSERT INTO users (id, email, password) VALUES ('5', 'jsmith@example.com', ?)",
			hashPassword("password"))
	}
}

func (s *UserRoutesTestSuite) AfterTest(suite, testName string) {
	s.T().Log("AfterTest: ", testName)
	switch testName {
	case "TestCreateUser":
		s.conn.Exec("DELETE FROM users WHERE email = 'ash@example.com'")
	}
}

func (s *UserRoutesTestSuite) TearDownSuite() {
	// cleanup
	s.conn.Exec("DELETE FROM users")
	s.conn.Close()
}

func (s *UserRoutesTestSuite) TestGetUsers() {
	req := httptest.NewRequest("GET", "/api/v1/users", nil)

	var users []db.User
	s.checkReqStatus(req, fiber.StatusOK, &users)

	if len(users) != 3 {
		s.T().Errorf("Expected 3 users but got %d", len(users))
	}
}

func (s *UserRoutesTestSuite) TestGetUser() {
	req := httptest.NewRequest("GET", "/api/v1/users/1", nil)

	var user db.User
	s.checkReqStatus(req, fiber.StatusOK, &user)

	s.Equal("John Doe", user.Name.String)
	s.Equal("1", user.ID)
}

func (s *UserRoutesTestSuite) TestGetUserNotFound() {
	req := httptest.NewRequest("GET", "/api/v1/users/100", nil)

	var user db.User
	s.checkReqStatus(req, fiber.StatusNotFound, &user)

	s.Equal("", user.ID)
}

func (s *UserRoutesTestSuite) TestCreateUser() {
	requestBody := []byte(`{"name": "Ashwin", "email": "ash@example.com", "password": "password"}`)
	req := httptest.NewRequest("POST", "/api/v1/users", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	var user db.User
	s.checkReqStatus(req, fiber.StatusCreated, &user)

	s.NotEmpty(user.ID, "ID should not be empty")
	s.Equal("Ashwin", user.Name.String)
	s.Equal("ash@example.com", user.Email)
	s.Empty(user.Password)
	s.NotEmpty(user.CreatedAt)
	s.NotEmpty(user.UpdatedAt)
}

func (s *UserRoutesTestSuite) TestCreateUserWithInvalidBody() {
	requestBody := []byte(`{"name": "Ashwin", "email": "ash@example.com", "password": "pass"}`)
	req := httptest.NewRequest("POST", "/api/v1/users", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	var errors []*utils.ErrorResponse
	s.checkReqStatus(req, fiber.StatusBadRequest, &errors)

	s.Equal(1, len(errors))
	s.Contains(strings.ToLower(errors[0].FailedField), "password")
}

func (s *UserRoutesTestSuite) TestUpdateUserName() {
	requestBody := []byte(`{"name": "Ashwin S"}`)
	req := httptest.NewRequest("PATCH", "/api/v1/users/1", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	var user db.User
	s.checkReqStatus(req, fiber.StatusOK, &user)

	s.Equal("Ashwin S", user.Name.String)
	s.Equal("1", user.ID)
}

func (s *UserRoutesTestSuite) TestDeleteUser() {
	req := httptest.NewRequest("DELETE", "/api/v1/users/5", nil)

	s.checkReqStatus(req, fiber.StatusNoContent, nil)

	var user db.User
	req = httptest.NewRequest("GET", "/api/v1/users/5", nil)

	s.checkReqStatus(req, fiber.StatusNotFound, &user)

	s.Equal("", user.ID)
	s.Equal("", user.Name.String)
}

func TestUserRoutes(t *testing.T) {
	suite.Run(t, new(UserRoutesTestSuite))
}

func (s *UserRoutesTestSuite) checkReqStatus(req *http.Request, expectedStatus int, out interface{}) {
	s.T().Helper()
	resp, err := s.app.Test(req, -1)
	s.NoError(err)

	s.Equal(expectedStatus, resp.StatusCode)

	if out != nil {
		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		s.T().Log(string(body))
		s.NoError(err)

		err = json.Unmarshal(body, &out)
		s.NoError(err)
	}
}
