package database

import (
	"context"
	"users-service/model"

	"github.com/google/uuid"
)

type UserDatabaseInterface interface {
	Create(ctx context.Context, user *model.User) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.User, error)
	List(ctx context.Context, page int, limit int, nameFilter string) ([]model.User, error)
}

type UserDatabase struct {
	pg *PostgresDB
}

func NewUserDatabase(db *PostgresDB) UserDatabaseInterface {
	return &UserDatabase{pg: db}
}

// Create a new user in the database
func (r *UserDatabase) Create(ctx context.Context, user *model.User) error {
	return r.pg.db.WithContext(ctx).Create(user).Error
}

// Get a user by ID
func (r *UserDatabase) GetByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	var user model.User
	err := r.pg.db.WithContext(ctx).First(&user, "id = ?", id).Error
	return &user, err
}

// List users with pagination and filter
func (r *UserDatabase) List(ctx context.Context, page int, limit int, nameFilter string) ([]model.User, error) {
	var users []model.User
	query := r.pg.db

	// Pagination
	if page == 0 {
		page = 1
	}
	offset := (page - 1) * limit

	// Filter on first name or last name with regular expression
	if nameFilter != "" {
		query = query.Where("lower(first_name) ~ ? OR lower(last_name) ~ ?",
			nameFilter, nameFilter)
	}

	// If limit is 0, return all users
	if limit == 0 {
		err := query.WithContext(ctx).Offset(offset).Find(&users).Error
		return users, err
	} else {
		err := query.WithContext(ctx).Offset(offset).Limit(limit).Find(&users).Error
		return users, err
	}

}
