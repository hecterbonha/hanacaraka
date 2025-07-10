package services

import (
	"errors"
	"testing"

	"com.hanacaraka/domain/entities"
)

// MockUserRepository is a mock implementation of UserRepository for testing
type MockUserRepository struct {
	users map[string]*entities.User
	// Control behavior for testing
	shouldFailGetAll    bool
	shouldFailGetByID   bool
	shouldFailCreate    bool
	shouldFailUpdate    bool
	shouldFailDelete    bool
}

func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{
		users: make(map[string]*entities.User),
	}
}

func (m *MockUserRepository) GetAll() ([]*entities.User, error) {
	if m.shouldFailGetAll {
		return nil, errors.New("repository error")
	}

	users := make([]*entities.User, 0, len(m.users))
	for _, user := range m.users {
		users = append(users, user)
	}
	return users, nil
}

func (m *MockUserRepository) GetByID(id string) (*entities.User, error) {
	if m.shouldFailGetByID {
		return nil, errors.New("repository error")
	}

	user, exists := m.users[id]
	if !exists {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (m *MockUserRepository) Create(user *entities.User) (*entities.User, error) {
	if m.shouldFailCreate {
		return nil, errors.New("repository error")
	}

	m.users[user.ID] = user
	return user, nil
}

func (m *MockUserRepository) Update(user *entities.User) (*entities.User, error) {
	if m.shouldFailUpdate {
		return nil, errors.New("repository error")
	}

	if _, exists := m.users[user.ID]; !exists {
		return nil, errors.New("user not found")
	}

	m.users[user.ID] = user
	return user, nil
}

func (m *MockUserRepository) Delete(id string) error {
	if m.shouldFailDelete {
		return errors.New("repository error")
	}

	if _, exists := m.users[id]; !exists {
		return errors.New("user not found")
	}

	delete(m.users, id)
	return nil
}



func TestNewUserService(t *testing.T) {
	mockRepo := NewMockUserRepository()
	service := NewUserService(mockRepo)

	if service == nil {
		t.Error("Expected service to be created, got nil")
	}

	if service.userRepo != mockRepo {
		t.Error("Expected service to have the provided repository")
	}
}

func TestUserService_GetAllUsers(t *testing.T) {
	tests := []struct {
		name           string
		setupRepo      func(*MockUserRepository)
		expectedError  bool
		expectedLength int
	}{
		{
			name: "Successfully get all users",
			setupRepo: func(repo *MockUserRepository) {
				user1 := entities.NewUser("John Doe", "john@example.com")
				user2 := entities.NewUser("Jane Smith", "jane@example.com")
				repo.users[user1.ID] = user1
				repo.users[user2.ID] = user2
			},
			expectedError:  false,
			expectedLength: 2,
		},
		{
			name: "Get all users from empty repository",
			setupRepo: func(repo *MockUserRepository) {
				// No users
			},
			expectedError:  false,
			expectedLength: 0,
		},
		{
			name: "Repository error",
			setupRepo: func(repo *MockUserRepository) {
				repo.shouldFailGetAll = true
			},
			expectedError:  true,
			expectedLength: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := NewMockUserRepository()
			tt.setupRepo(mockRepo)
			service := NewUserService(mockRepo)

			users, err := service.GetAllUsers()

			if tt.expectedError && err == nil {
				t.Error("Expected error, got nil")
			}
			if !tt.expectedError && err != nil {
				t.Errorf("Expected no error, got %v", err)
			}
			if len(users) != tt.expectedLength {
				t.Errorf("Expected %d users, got %d", tt.expectedLength, len(users))
			}
		})
	}
}

func TestUserService_GetUserByID(t *testing.T) {
	testID := "123e4567-e89b-12d3-a456-426614174000"
	testUser := entities.NewUserWithID(testID, "John Doe", "john@example.com")

	tests := []struct {
		name          string
		id            string
		setupRepo     func(*MockUserRepository)
		expectedError string
		expectedUser  *entities.User
	}{
		{
			name: "Successfully get user by ID",
			id:   testID,
			setupRepo: func(repo *MockUserRepository) {
				repo.users[testID] = testUser
			},
			expectedError: "",
			expectedUser:  testUser,
		},
		{
			name: "Invalid user ID - empty",
			id:   "",
			setupRepo: func(repo *MockUserRepository) {
				// No setup needed
			},
			expectedError: "invalid user ID",
			expectedUser:  nil,
		},
		{
			name: "User not found",
			id:   "nonexistent-id",
			setupRepo: func(repo *MockUserRepository) {
				// No users added
			},
			expectedError: "user not found",
			expectedUser:  nil,
		},
		{
			name: "Repository error",
			id:   testID,
			setupRepo: func(repo *MockUserRepository) {
				repo.shouldFailGetByID = true
			},
			expectedError: "repository error",
			expectedUser:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := NewMockUserRepository()
			tt.setupRepo(mockRepo)
			service := NewUserService(mockRepo)

			user, err := service.GetUserByID(tt.id)

			if tt.expectedError != "" {
				if err == nil {
					t.Errorf("Expected error '%s', got nil", tt.expectedError)
				} else if err.Error() != tt.expectedError {
					t.Errorf("Expected error '%s', got '%s'", tt.expectedError, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}
			}

			if tt.expectedUser != nil && user != nil {
				if user.ID != tt.expectedUser.ID || user.Name != tt.expectedUser.Name || user.Email != tt.expectedUser.Email {
					t.Errorf("Expected user %+v, got %+v", tt.expectedUser, user)
				}
			} else if tt.expectedUser != user {
				t.Errorf("Expected user %+v, got %+v", tt.expectedUser, user)
			}
		})
	}
}

func TestUserService_CreateUser(t *testing.T) {
	tests := []struct {
		name          string
		userName      string
		email         string
		setupRepo     func(*MockUserRepository)
		expectedError string
		expectUser    bool
	}{
		{
			name:          "Successfully create user",
			userName:      "John Doe",
			email:         "john@example.com",
			setupRepo:     func(repo *MockUserRepository) {},
			expectedError: "",
			expectUser:    true,
		},
		{
			name:          "Create user with empty name",
			userName:      "",
			email:         "john@example.com",
			setupRepo:     func(repo *MockUserRepository) {},
			expectedError: "invalid user data: name and email are required",
			expectUser:    false,
		},
		{
			name:          "Create user with empty email",
			userName:      "John Doe",
			email:         "",
			setupRepo:     func(repo *MockUserRepository) {},
			expectedError: "invalid user data: name and email are required",
			expectUser:    false,
		},
		{
			name:          "Create user with both empty name and email",
			userName:      "",
			email:         "",
			setupRepo:     func(repo *MockUserRepository) {},
			expectedError: "invalid user data: name and email are required",
			expectUser:    false,
		},
		{
			name:      "Repository error during create",
			userName:  "John Doe",
			email:     "john@example.com",
			setupRepo: func(repo *MockUserRepository) {
				repo.shouldFailCreate = true
			},
			expectedError: "repository error",
			expectUser:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := NewMockUserRepository()
			tt.setupRepo(mockRepo)
			service := NewUserService(mockRepo)

			user, err := service.CreateUser(tt.userName, tt.email)

			if tt.expectedError != "" {
				if err == nil {
					t.Errorf("Expected error '%s', got nil", tt.expectedError)
				} else if err.Error() != tt.expectedError {
					t.Errorf("Expected error '%s', got '%s'", tt.expectedError, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}
			}

			if tt.expectUser {
				if user == nil {
					t.Error("Expected user to be created, got nil")
				} else {
					if user.Name != tt.userName {
						t.Errorf("Expected user name '%s', got '%s'", tt.userName, user.Name)
					}
					if user.Email != tt.email {
						t.Errorf("Expected user email '%s', got '%s'", tt.email, user.Email)
					}
					// Check that ID is a valid UUID
					if user.ID == "" {
						t.Error("Expected user ID to be set")
					}
				}
			} else {
				if user != nil {
					t.Errorf("Expected no user to be created, got %+v", user)
				}
			}
		})
	}
}

func TestUserService_UpdateUser(t *testing.T) {
	testID := "123e4567-e89b-12d3-a456-426614174000"
	originalUser := entities.NewUserWithID(testID, "Original Name", "original@example.com")

	tests := []struct {
		name          string
		id            string
		userName      string
		email         string
		setupRepo     func(*MockUserRepository)
		expectedError string
		expectUser    bool
	}{
		{
			name:     "Successfully update user",
			id:       testID,
			userName: "Updated Name",
			email:    "updated@example.com",
			setupRepo: func(repo *MockUserRepository) {
				repo.users[testID] = originalUser
			},
			expectedError: "",
			expectUser:    true,
		},
		{
			name:     "Update user with empty name",
			id:       testID,
			userName: "",
			email:    "updated@example.com",
			setupRepo: func(repo *MockUserRepository) {
				repo.users[testID] = originalUser
			},
			expectedError: "invalid user data: name and email are required",
			expectUser:    false,
		},
		{
			name:     "Update user with empty email",
			id:       testID,
			userName: "Updated Name",
			email:    "",
			setupRepo: func(repo *MockUserRepository) {
				repo.users[testID] = originalUser
			},
			expectedError: "invalid user data: name and email are required",
			expectUser:    false,
		},
		{
			name:          "Invalid user ID - empty",
			id:            "",
			userName:      "Updated Name",
			email:         "updated@example.com",
			setupRepo:     func(repo *MockUserRepository) {},
			expectedError: "invalid user ID",
			expectUser:    false,
		},
		{
			name:          "User not found",
			id:            "nonexistent-id",
			userName:      "Updated Name",
			email:         "updated@example.com",
			setupRepo:     func(repo *MockUserRepository) {},
			expectedError: "user not found",
			expectUser:    false,
		},
		{
			name:     "Repository error during GetByID",
			id:       testID,
			userName: "Updated Name",
			email:    "updated@example.com",
			setupRepo: func(repo *MockUserRepository) {
				repo.shouldFailGetByID = true
			},
			expectedError: "repository error",
			expectUser:    false,
		},
		{
			name:     "Repository error during Update",
			id:       testID,
			userName: "Updated Name",
			email:    "updated@example.com",
			setupRepo: func(repo *MockUserRepository) {
				repo.users[testID] = originalUser
				repo.shouldFailUpdate = true
			},
			expectedError: "repository error",
			expectUser:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := NewMockUserRepository()
			tt.setupRepo(mockRepo)
			service := NewUserService(mockRepo)

			user, err := service.UpdateUser(tt.id, tt.userName, tt.email)

			if tt.expectedError != "" {
				if err == nil {
					t.Errorf("Expected error '%s', got nil", tt.expectedError)
				} else if err.Error() != tt.expectedError {
					t.Errorf("Expected error '%s', got '%s'", tt.expectedError, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}
			}

			if tt.expectUser {
				if user == nil {
					t.Error("Expected user to be updated, got nil")
				} else {
					if user.ID != tt.id {
						t.Errorf("Expected user ID %s, got %s", tt.id, user.ID)
					}
					if user.Name != tt.userName {
						t.Errorf("Expected user name '%s', got '%s'", tt.userName, user.Name)
					}
					if user.Email != tt.email {
						t.Errorf("Expected user email '%s', got '%s'", tt.email, user.Email)
					}
				}
			} else {
				if user != nil {
					t.Errorf("Expected no user to be returned, got %+v", user)
				}
			}
		})
	}
}

func TestUserService_DeleteUser(t *testing.T) {
	testID := "123e4567-e89b-12d3-a456-426614174000"
	testUser := entities.NewUserWithID(testID, "John Doe", "john@example.com")

	tests := []struct {
		name          string
		id            string
		setupRepo     func(*MockUserRepository)
		expectedError string
	}{
		{
			name: "Successfully delete user",
			id:   testID,
			setupRepo: func(repo *MockUserRepository) {
				repo.users[testID] = testUser
			},
			expectedError: "",
		},
		{
			name:          "Invalid user ID - empty",
			id:            "",
			setupRepo:     func(repo *MockUserRepository) {},
			expectedError: "invalid user ID",
		},
		{
			name:          "User not found",
			id:            "nonexistent-id",
			setupRepo:     func(repo *MockUserRepository) {},
			expectedError: "user not found",
		},
		{
			name: "Repository error during get",
			id:   testID,
			setupRepo: func(repo *MockUserRepository) {
				repo.shouldFailGetByID = true
			},
			expectedError: "repository error",
		},
		{
			name: "Repository error during delete",
			id:   testID,
			setupRepo: func(repo *MockUserRepository) {
				repo.users[testID] = testUser
				repo.shouldFailDelete = true
			},
			expectedError: "repository error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := NewMockUserRepository()
			tt.setupRepo(mockRepo)
			service := NewUserService(mockRepo)

			err := service.DeleteUser(tt.id)

			if tt.expectedError != "" {
				if err == nil {
					t.Errorf("Expected error '%s', got nil", tt.expectedError)
				} else if err.Error() != tt.expectedError {
					t.Errorf("Expected error '%s', got '%s'", tt.expectedError, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}
			}
		})
	}
}
