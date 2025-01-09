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

func (r *UserDatabase) Create(ctx context.Context, user *model.User) error {
	return r.pg.db.WithContext(ctx).Create(user).Error
}

func (r *UserDatabase) GetByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	var user model.User
	err := r.pg.db.WithContext(ctx).First(&user, "id = ?", id).Error
	return &user, err
}

func (r *UserDatabase) List(ctx context.Context, page int, limit int, nameFilter string) ([]model.User, error) {
	var users []model.User
	query := r.pg.db

	if page == 0 {
		page = 1
	}
	offset := (page - 1) * limit

	if nameFilter != "" {
		query = query.Where("lower(first_name) ~ ? OR lower(last_name) ~ ?",
			nameFilter, nameFilter)
	}

	if limit == 0 {
		err := query.WithContext(ctx).Offset(offset).Find(&users).Error
		return users, err
	} else {
		err := query.WithContext(ctx).Offset(offset).Limit(limit).Find(&users).Error
		return users, err
	}

}
