package service

import (
	"testing"
	"time"

	"users-service/database"
	"users-service/model"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type mockDatabase struct {
	users map[uuid.UUID]*model.User
}

func newMockDatabase() database.UserDatabaseInterface {
	return &mockDatabase{
		users: make(map[uuid.UUID]*model.User),
	}
}

func (m *mockDatabase) Create(user *model.User) error {
	m.users[user.ID] = user
	return nil
}

func (m *mockDatabase) GetByID(id uuid.UUID) (*model.User, error) {
	if user, exists := m.users[id]; exists {
		return user, nil
	}
	return nil, nil
}

func (m *mockDatabase) List(limit int, nameFilter string) ([]model.User, error) {
	var users []model.User
	for _, user := range m.users {
		if nameFilter == "" || user.FirstName == nameFilter || user.LastName == nameFilter {
			users = append(users, *user)
		}
	}
	if limit > 0 && len(users) > limit {
		users = users[:limit]
	}
	return users, nil
}

type mockRandomUserClient struct {
	response *RandomUserResponse
	error    error
}

func newMockRandomUserClient() *mockRandomUserClient {
	return &mockRandomUserClient{
		response: &RandomUserResponse{
			Results: []struct {
				Name struct {
					First string `json:"first"`
					Last  string `json:"last"`
				} `json:"name"`
				Gender string `json:"gender"`
				Email  string `json:"email"`
				Phone  string `json:"phone"`
			}{
				{
					Name: struct {
						First string `json:"first"`
						Last  string `json:"last"`
					}{
						First: "John",
						Last:  "Doe",
					},
					Gender: "male",
					Email:  "john.doe@example.com",
					Phone:  "123-456-7890",
				},
			},
		},
	}
}

func (m *mockRandomUserClient) GetRandomUsers(count int, gender string) (*RandomUserResponse, error) {
	if m.error != nil {
		return nil, m.error
	}
	return m.response, nil
}

func TestValidateGender(t *testing.T) {
	tests := []struct {
		name    string
		gender  string
		wantErr bool
	}{
		{
			name:    "Valid male gender",
			gender:  "male",
			wantErr: false,
		},
		{
			name:    "Valid female gender",
			gender:  "female",
			wantErr: false,
		},
		{
			name:    "Valid any gender",
			gender:  "any",
			wantErr: false,
		},
		{
			name:    "Invalid gender",
			gender:  "invalid",
			wantErr: true,
		},
		{
			name:    "Empty gender",
			gender:  "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateGender(tt.gender)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestUserService_CreateUsers(t *testing.T) {
	tests := []struct {
		name    string
		count   int
		gender  string
		wantErr bool
		setup   func(*mockRandomUserClient)
		check   func(*testing.T, *model.CreateUsersResponse)
	}{
		{
			name:    "Successful creation with valid gender",
			count:   1,
			gender:  "male",
			wantErr: false,
			setup:   func(m *mockRandomUserClient) {},
			check: func(t *testing.T, resp *model.CreateUsersResponse) {
				assert.Equal(t, 1, resp.RequestedCount)
				assert.Equal(t, "male", resp.Gender)
			},
		},
		{
			name:    "Invalid gender",
			count:   1,
			gender:  "invalid",
			wantErr: true,
			setup:   func(m *mockRandomUserClient) {},
			check: func(t *testing.T, resp *model.CreateUsersResponse) {
				assert.Nil(t, resp)
			},
		},
		{
			name:    "API error",
			count:   1,
			gender:  "male",
			wantErr: true,
			setup: func(m *mockRandomUserClient) {
				m.error = assert.AnError
			},
			check: func(t *testing.T, resp *model.CreateUsersResponse) {
				assert.Nil(t, resp)
			},
		},
		{
			name:    "Zero count",
			count:   0,
			gender:  "male",
			wantErr: false,
			setup:   func(m *mockRandomUserClient) {},
			check: func(t *testing.T, resp *model.CreateUsersResponse) {
				assert.Equal(t, 0, resp.RequestedCount)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := newMockRandomUserClient()
			tt.setup(mockClient)

			service := NewUserService(newMockDatabase(), mockClient)
			response, err := service.CreateUsers(tt.count, tt.gender)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			tt.check(t, response)
		})
	}
}

func TestUserService_GetUser(t *testing.T) {
	db := newMockDatabase()
	service := NewUserService(db, newMockRandomUserClient())

	// Test user
	testUser := &model.User{
		ID:        uuid.New(),
		FirstName: "Test",
		LastName:  "User",
		Gender:    "male",
		Email:     "test@example.com",
		Phone:     "123-456-7890",
		CreatedAt: time.Now(),
	}
	db.Create(testUser)

	t.Run("Get existing user", func(t *testing.T) {
		user, err := service.GetUser(testUser.ID)
		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, testUser.ID, user.ID)
		assert.Equal(t, testUser.FirstName, user.FirstName)
	})

	t.Run("Get non-existent user", func(t *testing.T) {
		user, err := service.GetUser(uuid.New())
		assert.NoError(t, err)
		assert.Nil(t, user)
	})
}

func TestUserService_ListUsers(t *testing.T) {
	db := newMockDatabase()
	service := NewUserService(db, newMockRandomUserClient())

	// Create test users
	testUsers := []struct {
		firstName string
		lastName  string
	}{
		{"Alice", "Smith"},
		{"Bob", "Johnson"},
		{"Charlie", "Brown"},
	}

	for _, u := range testUsers {
		user := &model.User{
			ID:        uuid.New(),
			FirstName: u.firstName,
			LastName:  u.lastName,
			CreatedAt: time.Now(),
		}
		db.Create(user)
	}

	t.Run("List all users", func(t *testing.T) {
		users, err := service.ListUsers(0, "")
		assert.NoError(t, err)
		assert.Len(t, users, len(testUsers))
	})

	t.Run("List with limit", func(t *testing.T) {
		users, err := service.ListUsers(2, "")
		assert.NoError(t, err)
		assert.Len(t, users, 2)
	})
}
