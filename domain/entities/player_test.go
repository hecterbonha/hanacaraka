package entities

import (
	"testing"

	"github.com/google/uuid"
)

func TestNewPlayer(t *testing.T) {
	tests := []struct {
		name     string
		userName string
		email    string
	}{
		{
			name:     "Valid user creation",
			userName: "John Doe",
			email:    "john@example.com",
		},
		{
			name:     "User with empty name",
			userName: "",
			email:    "test@example.com",
		},
		{
			name:     "User with empty email",
			userName: "Test User",
			email:    "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := NewUser(tt.userName, tt.email)

			// Check that ID is a valid UUID
			if _, err := uuid.Parse(user.ID); err != nil {
				t.Errorf("Expected valid UUID for ID, got %s", user.ID)
			}
			if user.Name != tt.userName {
				t.Errorf("Expected Name %s, got %s", tt.userName, user.Name)
			}
			if user.Email != tt.email {
				t.Errorf("Expected Email %s, got %s", tt.email, user.Email)
			}
		})
	}
}

func TestNewPlayerWithID(t *testing.T) {
	tests := []struct {
		name     string
		id       string
		userName string
		email    string
	}{
		{
			name:     "Valid user creation with specific ID",
			id:       "123e4567-e89b-12d3-a456-426614174000",
			userName: "John Doe",
			email:    "john@example.com",
		},
		{
			name:     "User with empty ID",
			id:       "",
			userName: "Jane Smith",
			email:    "jane@example.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := NewUserWithID(tt.id, tt.userName, tt.email)

			if user.ID != tt.id {
				t.Errorf("Expected ID %s, got %s", tt.id, user.ID)
			}
			if user.Name != tt.userName {
				t.Errorf("Expected Name %s, got %s", tt.userName, user.Name)
			}
			if user.Email != tt.email {
				t.Errorf("Expected Email %s, got %s", tt.email, user.Email)
			}
		})
	}
}

func TestPlayer_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		user     *User
		expected bool
	}{
		{
			name:     "Valid user with name and email",
			user:     &User{ID: "123e4567-e89b-12d3-a456-426614174000", Name: "John Doe", Email: "john@example.com"},
			expected: true,
		},
		{
			name:     "Invalid user with empty name",
			user:     &User{ID: "123e4567-e89b-12d3-a456-426614174000", Name: "", Email: "john@example.com"},
			expected: false,
		},
		{
			name:     "Invalid user with empty email",
			user:     &User{ID: "123e4567-e89b-12d3-a456-426614174000", Name: "John Doe", Email: ""},
			expected: false,
		},
		{
			name:     "Invalid user with both empty name and email",
			user:     &User{ID: "123e4567-e89b-12d3-a456-426614174000", Name: "", Email: ""},
			expected: false,
		},
		{
			name:     "Valid user with whitespace in name and email",
			user:     &User{ID: "123e4567-e89b-12d3-a456-426614174000", Name: " John Doe ", Email: " john@example.com "},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.user.IsValid()
			if result != tt.expected {
				t.Errorf("Expected IsValid() to return %v, got %v", tt.expected, result)
			}
		})
	}
}
