package mocks

import (
	"time"

	"snippetbox.stanley.net/internal/models"
)

const (
	defaultUserId       = 1
	defaultUserName     = "Alice Bob"
	defaultUserEmail    = "alice@example.com"
	defaultUserPassword = "pa$$word"
)

type UserModel struct{}

func (m *UserModel) Get(id int) (*models.User, error) {
	switch id {
	case 1:
		user := &models.User{
			Id: defaultUserId,	
			Name: defaultUserName,
			Email: defaultUserEmail,
			Created: time.Now().UTC(),
		}
		return user, nil
	default:
		return nil, models.ErrNoRecord
	}
}

func (m *UserModel) Insert(name, email, password string) error {
	switch email {
	case defaultUserEmail:
		return models.ErrDuplicateEmail
	default:
		return nil
	}
}

func (m *UserModel) Authenticate(email, password string) (int, error) {
	if email == defaultUserEmail && password == defaultUserPassword {
		return defaultUserId, nil
	}

	return 0, models.ErrInvalidCredentials
}

func (m *UserModel) Exists(id int) (bool, error) {
	switch id {
	case defaultUserId:
		return true, nil
	default:
		return false, nil
	}
}

func (m *UserModel) UpdatePassword(id int, currentPassword, newPassword string) error {
	if id != defaultUserId {
		return models.ErrNoRecord
	}

	if currentPassword != defaultUserPassword {
		return models.ErrInvalidCredentials
	}

	return nil
}
