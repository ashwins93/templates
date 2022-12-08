package db

import (
	"encoding/gob"
	"fmt"

	"github.com/ashwins93/fiber-badger/utils"
	"github.com/dgraph-io/badger/v3"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Username     string `json:"username"`
	PasswordHash string `json:"-"`
	FirstName    string `json:"firstName"`
	LastName     string `json:"lastName"`
}

type CreateUserParams struct {
	Username  string `json:"username" validate:"required,min=6,max=25"`
	Password  string `json:"password" validate:"required,min=8,max=15"`
	FirstName string `json:"firstName" validate:"min=2,alpha"`
	LastName  string `json:"lastName" validate:"min=2,alpha"`
}

func init() {
	gob.Register(User{})
}

func (q *Queries) CreateNewUser(data *CreateUserParams) (*User, error) {
	var user *User
	err := q.db.Update(func(txn *badger.Txn) error {
		key := []byte(fmt.Sprintf("user/%s", data.Username))
		_, err := txn.Get(key)

		if err == nil {
			return fiber.NewError(fiber.StatusBadRequest, "Username already taken")
		} else if err != nil && err != badger.ErrKeyNotFound {
			return err
		}

		passwordHash, err := bcrypt.GenerateFromPassword([]byte(data.Password), 10)
		if err != nil {
			return err
		}

		user = &User{
			Username:     data.Username,
			PasswordHash: string(passwordHash),
			FirstName:    data.FirstName,
			LastName:     data.LastName,
		}

		value, err := utils.MarshalStruct(user)
		if err != nil {
			return err
		}

		err = txn.Set(key, value)

		return err
	})

	return user, err
}

func (q *Queries) GetAllUsers() ([]*User, error) {
	users := make([]*User, 0)
	prefix := []byte("user/")
	err := q.db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.IteratorOptions{
			PrefetchSize: 10,
			Prefix:       prefix,
		})
		defer it.Close()

		for it.Seek(prefix); it.Valid(); it.Next() {
			item := it.Item()
			err := item.Value(func(v []byte) error {
				user, err := utils.UnmarshalStruct[User](v)
				users = append(users, &user)
				return err
			})
			if err != nil {
				return err
			}
		}

		return nil
	})

	return users, err
}
