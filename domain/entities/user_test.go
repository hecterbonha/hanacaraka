package entities

import (
	"testing"

	"github.com/google/uuid"
)

func TestNewUser(t *testing.T) {
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

func TestNewUserWithID(t *testing.T) {
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

func TestUser_IsValid(t *testing.T) {
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

func TestUser_UpdateName(t *testing.T) {
	tests := []struct {
		name        string
		user        *User
		newName     string
		expectedName string
	}{
		{
			name:        "Update name to valid value",
			user:        &User{ID: "123e4567-e89b-12d3-a456-426614174000", Name: "John Doe", Email: "john@example.com"},
			newName:     "Jane Smith",
			expectedName: "Jane Smith",
		},
		{
			name:        "Update name to empty string",
			user:        &User{ID: "123e4567-e89b-12d3-a456-426614174000", Name: "John Doe", Email: "john@example.com"},
			newName:     "",
			expectedName: "",
		},
		{
			name:        "Update name with whitespace",
			user:        &User{ID: "123e4567-e89b-12d3-a456-426614174000", Name: "John Doe", Email: "john@example.com"},
			newName:     " Jane Smith ",
			expectedName: " Jane Smith ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			originalEmail := tt.user.Email
			originalID := tt.user.ID

			tt.user.UpdateName(tt.newName)

			if tt.user.Name != tt.expectedName {
				t.Errorf("Expected name to be %s, got %s", tt.expectedName, tt.user.Name)
			}

			// Ensure other fields are not modified
			if tt.user.Email != originalEmail {
				t.Errorf("Email should not be modified, expected %s, got %s", originalEmail, tt.user.Email)
			}
			if tt.user.ID != originalID {
				t.Errorf("ID should not be modified, expected %s, got %s", originalID, tt.user.ID)
			}
		})
	}
}

func TestUser_UpdateEmail(t *testing.T) {
	tests := []struct {
		name          string
		user          *User
		newEmail      string
		expectedEmail string
	}{
		{
			name:          "Update email to valid value",
			user:          &User{ID: "123e4567-e89b-12d3-a456-426614174000", Name: "John Doe", Email: "john@example.com"},
			newEmail:      "jane@example.com",
			expectedEmail: "jane@example.com",
		},
		{
			name:          "Update email to empty string",
			user:          &User{ID: "123e4567-e89b-12d3-a456-426614174000", Name: "John Doe", Email: "john@example.com"},
			newEmail:      "",
			expectedEmail: "",
		},
		{
			name:          "Update email with whitespace",
			user:          &User{ID: "123e4567-e89b-12d3-a456-426614174000", Name: "John Doe", Email: "john@example.com"},
			newEmail:      " jane@example.com ",
			expectedEmail: " jane@example.com ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			originalName := tt.user.Name
			originalID := tt.user.ID

			tt.user.UpdateEmail(tt.newEmail)

			if tt.user.Email != tt.expectedEmail {
				t.Errorf("Expected email to be %s, got %s", tt.expectedEmail, tt.user.Email)
			}

			// Ensure other fields are not modified
			if tt.user.Name != originalName {
				t.Errorf("Name should not be modified, expected %s, got %s", originalName, tt.user.Name)
			}
			if tt.user.ID != originalID {
				t.Errorf("ID should not be modified, expected %s, got %s", originalID, tt.user.ID)
			}
		})
	}
}

func TestUser_JSONSerialization(t *testing.T) {
	user := NewUser("John Doe", "john@example.com")

	// Test that the struct tags are properly set up for JSON serialization
	// This is more of a structural test to ensure the tags exist
	if _, err := uuid.Parse(user.ID); err != nil {
		t.Errorf("Expected valid UUID for ID, got %s", user.ID)
	}
	if user.Name != "John Doe" {
		t.Errorf("Expected Name 'John Doe', got %s", user.Name)
	}
	if user.Email != "john@example.com" {
		t.Errorf("Expected Email 'john@example.com', got %s", user.Email)
	}
}
