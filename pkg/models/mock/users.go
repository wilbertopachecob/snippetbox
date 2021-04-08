package mock

import (
	"time"
	"wilbertopachecob/snippetbox/pkg/models"
)

type UserModel struct{}

var mockUser = &models.User{
	ID:      1,
	Name:    "Admin",
	Email:   "admin@gmail.com",
	Created: time.Now(),
}

func (m *UserModel) Insert(name, email, password string) error {
	switch email {
	case "admin@gmail.com":
		return nil
	default:
		return models.ErrDuplicateEmail
	}
}

func (m *UserModel) Authenticate(email, password string) (int, error) {
	switch email {
	case "admin@gmail.com":
		return 1, nil
	default:
		return 0, models.ErrInvalidCredentials
	}
}

func (m *UserModel) Get(ID int) (*models.User, error) {
	switch ID {
	case 1:
		return mockUser, nil
	default:
		return nil, models.ErrNoRecord
	}
}
