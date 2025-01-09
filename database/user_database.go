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
	pg *PostgresDB
}

func NewUserDatabase(db *PostgresDB) UserDatabaseInterface {
	return &UserDatabase{pg: db}
}

func (r *UserDatabase) Create(user *model.User) error {
	return r.pg.db.Create(user).Error
}

func (r *UserDatabase) GetByID(id uuid.UUID) (*model.User, error) {
	var user model.User
	err := r.pg.db.First(&user, "id = ?", id).Error
	return &user, err
}

func (r *UserDatabase) List(limit int, nameFilter string) ([]model.User, error) {
	var users []model.User
	query := r.pg.db

	if limit > 0 {
		query = query.Limit(limit)
	}

	if nameFilter != "" {
		query = query.Where("lower(first_name) ~ ? OR lower(last_name) ~ ?",
			nameFilter, nameFilter)
	}

	err := query.Find(&users).Error
	return users, err
}
