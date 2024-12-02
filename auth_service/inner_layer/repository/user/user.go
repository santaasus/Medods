package user

import (
	"Medods/auth_service/inner_layer/db"
	domain "Medods/auth_service/inner_layer/domain"
	"errors"

	domainErrors "github.com/santaasus/errors-handler"

	_ "github.com/lib/pq"
)

type IRepository interface {
	GetUserByGuid(guid string) (*domain.User, error)
	CreateUser(newUser *domain.NewUser) (*domain.User, error)
	UpdateUser(params map[string]any, userId int) error
	DeleteUserByHash(hash string) error
}

type Repository struct {
}

func (Repository) GetUserByGuid(guid string) (*domain.User, error) {
	user, err := db.GetUserByGuid(guid)
	if err != nil {
		return nil, &domainErrors.AppError{
			Err:  errors.New("user does not exists"),
			Type: domainErrors.ValidationError,
		}
	}

	return user, nil
}

func (Repository) CreateUser(newUser *domain.NewUser) (*domain.User, error) {
	user, err := db.CreateUser(newUser)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (Repository) UpdateUser(params map[string]any, userId int) error {
	if len(params) == 0 {
		return &domainErrors.AppError{
			Err:  errors.New("UpdateUser: params are empty"),
			Type: domainErrors.ValidationError,
		}
	}

	_, err := db.UpdateUserByParams(params, userId)
	if err != nil {
		return err
	}

	return nil
}

func (Repository) DeleteUserByHash(hash string) error {
	err := db.DeleteUserByHash(hash)
	if err != nil {
		return &domainErrors.AppError{
			Err:  errors.New("user does not exists"),
			Type: domainErrors.ValidationError,
		}
	}

	return nil
}
