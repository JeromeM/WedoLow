package database

import (
	"users-service/model"

	"github.com/google/uuid"
)

type UserDatabaseInterface interface {
	Create(user *model.User) error
	GetByID(id uuid.UUID) (*model.User, error)
	List(limit int, nameFilter string) ([]model.User, error)
}

type UserDatabase struct {
	db *PostgresDB
}

func NewUserDatabase(db *PostgresDB) UserDatabaseInterface {
	return &UserDatabase{db: db}
}

func (r *UserDatabase) Create(user *model.User) error {
	return r.db.db.Create(user).Error
}

func (r *UserDatabase) GetByID(id uuid.UUID) (*model.User, error) {
	var user model.User
	err := r.db.db.First(&user, "id = ?", id).Error
	return &user, err
}

func (r *UserDatabase) List(limit int, nameFilter string) ([]model.User, error) {
	var users []model.User
	query := r.db.db

	err := query.Find(&users).Error
	return users, err
}
