package service

import (
	"time"
	"users-service/database"
	"users-service/model"

	"github.com/google/uuid"
)

type RandomUserClientInterface interface {
	GetRandomUsers(count int, gender string) (*RandomUserResponse, error)
}

type UserService struct {
	db         database.UserDatabaseInterface
	randomUser RandomUserClientInterface
}

func NewUserService(db database.UserDatabaseInterface, randomUser RandomUserClientInterface) *UserService {
	return &UserService{
		db:         db,
		randomUser: randomUser,
	}
}

func (s *UserService) CreateUsers(count int, gender string) error {
	resp, err := s.randomUser.GetRandomUsers(count, gender)
	if err != nil {
		return err
	}

	for _, result := range resp.Results {
		user := &model.User{
			ID:        uuid.New(),
			FirstName: result.Name.First,
			LastName:  result.Name.Last,
			Gender:    result.Gender,
			Email:     result.Email,
			Phone:     result.Phone,
			CreatedAt: time.Now(),
		}
		if err := s.db.Create(user); err != nil {
			return err
		}
	}

	return nil
}

func (s *UserService) GetUser(id uuid.UUID) (*model.User, error) {
	return s.db.GetByID(id)
}

func (s *UserService) ListUsers(limit int, nameFilter string) ([]model.User, error) {
	return s.db.List(limit, nameFilter)
}
