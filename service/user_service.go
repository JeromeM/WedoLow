package service

import (
	"fmt"
	"time"
	"users-service/database"
	"users-service/model"

	"github.com/google/uuid"
)

const (
	GenderMale   = "male"
	GenderFemale = "female"
	GenderAny    = "any"
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

func ValidateGender(gender string) error {
	switch gender {
	case GenderMale, GenderFemale, GenderAny:
		return nil
	default:
		return fmt.Errorf("invalid gender : %s. Authorized values are : male, female, any", gender)
	}
}

func (s *UserService) CreateUsers(count int, gender string) (*model.CreateUsersResponse, error) {

	if err := ValidateGender(gender); err != nil {
		return nil, err
	}

	resp, err := s.randomUser.GetRandomUsers(count, gender)
	if err != nil {
		return nil, fmt.Errorf("Error with Random users public API : %w", err)
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
			return nil, fmt.Errorf("erreur lors de la création de l'utilisateur : %w", err)
		}
	}

	response := &model.CreateUsersResponse{
		RequestedCount: count,
		Gender:         gender,
	}

	return response, nil
}

func (s *UserService) GetUser(id uuid.UUID) (*model.User, error) {
	return s.db.GetByID(id)
}

func (s *UserService) ListUsers(limit int, nameFilter string) ([]model.User, error) {
	return s.db.List(limit, nameFilter)
}
