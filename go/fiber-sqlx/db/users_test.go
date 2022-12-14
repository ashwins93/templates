package db

import (
	"context"
	"testing"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/suite"
	"golang.org/x/crypto/bcrypt"
)

type UsersTestSuite struct {
	suite.Suite
	q    *Queries
	conn *sqlx.DB
}

func (s *UsersTestSuite) SetupSuite() {
	conn, err := sqlx.Connect("sqlite3", "file:test.db?_fk=1")
	if err != nil {
		panic(err)
	}

	s.q = NewDb(conn)
	s.conn = conn

	userList := []struct {
		id    string
		name  string
		email string
	}{
		{"1", "Jane Doe", "janedoe@example.com"},
		{"2", "John Smith", "johnsmith@example.com"},
		{"3", "Jane Smith", "janesmith@example.com"},
	}

	for _, user := range userList {
		var name NullString
		name.String = user.name
		name.Valid = true
		userParams := CreateUserParams{
			ID:       user.id,
			Email:    user.email,
			Password: hashPassword("password"),
			Name:     name,
		}
		s.insertUser(userParams)
	}
}

func (s *UsersTestSuite) BeforeTest(suite, testName string) {
	s.T().Log("BeforeTest: ", testName)
	switch testName {
	case "TestDeleteUser":
		s.insertUser(CreateUserParams{
			ID:       "5",
			Email:    "jsmith@example.com",
			Password: hashPassword("password"),
		})
	}
}

func (s *UsersTestSuite) TearDownSuite() {
	// cleanup
	s.q.db.ExecContext(context.Background(), "DELETE FROM users")
	s.conn.Close()
}

func (s *UsersTestSuite) TestCreateUser() {
	var name NullString
	name.String = "John Doe"
	name.Valid = true

	userParams := CreateUserParams{
		ID:       "4",
		Email:    "johndoe@example.com",
		Password: hashPassword("password"),
		Name:     name,
	}

	s.insertUser(userParams)
}

func (s *UsersTestSuite) TestGetUsers() {
	users, err := s.q.GetUsers(context.Background())
	s.NoError(err)
	s.Len(users, 4)
}

func (s *UsersTestSuite) TestGetUserByEmail() {
	email := "johndoe@example.com"
	user, err := s.q.GetUserByEmail(context.Background(), email)
	s.NoError(err)
	s.Equal(email, user.Email)
}

func (s *UsersTestSuite) TestGetUserByID() {
	id := "1"
	user, err := s.q.GetUserByID(context.Background(), id)
	s.NoError(err)
	s.Equal(id, user.ID)
}

func (s *UsersTestSuite) TestUpdateUser() {
	var name NullString
	name.String = "Jane Dot"
	name.Valid = true
	id := "1"

	var password NullString
	password.String = hashPassword("newpassword")
	password.Valid = true

	updateParams := UpdateUserParams{
		ID:       id,
		Name:     name,
		Password: password,
	}

	updatedUser, err := s.q.UpdateUser(context.Background(), updateParams)
	s.NoError(err)
	s.Equal(name.String, updatedUser.Name.String)
	s.Equal(password.String, updatedUser.Password)
	s.Equal("janedoe@example.com", updatedUser.Email, "Email is unchanged")
}

func (s *UsersTestSuite) TestDeleteUser() {
	id := "5"

	err := s.q.DeleteUser(context.Background(), id)
	s.NoError(err)

	user, err := s.q.GetUserByEmail(context.Background(), id)
	s.NoError(err)
	s.Equal("", user.Email)
}

func (s *UsersTestSuite) TestGetUserByEmailNotFound() {
	email := "nonexistent@example.com"
	user, err := s.q.GetUserByEmail(context.Background(), email)
	s.NoError(err)
	s.Equal("", user.Email)
}

func (s *UsersTestSuite) TestPartialUpdates() {
	var name NullString
	name.String = "Jane Smith"
	name.Valid = true
	id := "2"

	updateParams := UpdateUserParams{
		ID:       id,
		Name:     name,
		Password: NullString{},
	}

	updatedUser, err := s.q.UpdateUser(context.Background(), updateParams)
	s.NoError(err)
	s.Equal(name.String, updatedUser.Name.String)
	s.True(isPasswordSame("password", updatedUser.Password))
}

func TestUsers(t *testing.T) {
	suite.Run(t, new(UsersTestSuite))
}

func (s *UsersTestSuite) insertUser(userParams CreateUserParams) (*User, error) {
	s.T().Helper()
	user, err := s.q.CreateUser(context.Background(), userParams)
	s.NoError(err)

	s.Assert().Equal(userParams.Email, user.Email)
	s.Assert().Equal(userParams.Password, user.Password)
	s.Assert().Equal(userParams.Name.String, user.Name.String)
	s.Assert().NotEqual(0, user.CreatedAt)
	s.Assert().NotEqual(0, user.UpdatedAt)

	return user, err
}

func hashPassword(password string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	return string(hash)
}

func isPasswordSame(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
