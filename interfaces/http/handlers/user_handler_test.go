package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"com.hanacaraka/domain/entities"

	"github.com/gorilla/mux"
)

// MockUserService is a mock implementation of services.UserServiceInterface for testing
type MockUserService struct {
	users            map[string]*entities.User
	shouldFailGet    bool
	shouldFailCreate bool
	shouldFailUpdate bool
	shouldFailDelete bool
}

func NewMockUserService() *MockUserService {
	return &MockUserService{
		users: make(map[string]*entities.User),
	}
}

func (m *MockUserService) GetAllUsers() ([]*entities.User, error) {
	if m.shouldFailGet {
		return nil, errors.New("service error")
	}

	users := make([]*entities.User, 0, len(m.users))
	for _, user := range m.users {
		users = append(users, user)
	}
	return users, nil
}

func (m *MockUserService) GetUserByID(id string) (*entities.User, error) {
	if m.shouldFailGet {
		return nil, errors.New("service error")
	}

	user, exists := m.users[id]
	if !exists {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (m *MockUserService) CreateUser(name, email string) (*entities.User, error) {
	if m.shouldFailCreate {
		return nil, errors.New("service error")
	}

	if name == "" || email == "" {
		return nil, errors.New("invalid user data: name and email are required")
	}

	user := entities.NewUser(name, email)
	m.users[user.ID] = user
	return user, nil
}

func (m *MockUserService) UpdateUser(id string, name, email string) (*entities.User, error) {
	if m.shouldFailUpdate {
		return nil, errors.New("service error")
	}

	user, exists := m.users[id]
	if !exists {
		return nil, errors.New("user not found")
	}

	if name == "" || email == "" {
		return nil, errors.New("invalid user data: name and email are required")
	}

	user.Name = name
	user.Email = email
	return user, nil
}

func (m *MockUserService) DeleteUser(id string) error {
	if m.shouldFailDelete {
		return errors.New("service error")
	}

	_, exists := m.users[id]
	if !exists {
		return errors.New("user not found")
	}

	delete(m.users, id)
	return nil
}

func TestNewUserHandler(t *testing.T) {
	mockService := NewMockUserService()
	handler := NewUserHandler(mockService)

	if handler == nil {
		t.Error("Expected handler to be created, got nil")
	}

	// We can't directly compare interfaces, so we'll just check that it's not nil
	if handler.userService == nil {
		t.Error("Expected handler to have a service")
	}
}

func TestUserHandler_GetUsers(t *testing.T) {
	tests := []struct {
		name           string
		setupService   func(*MockUserService)
		expectedStatus int
		expectedUsers  int
	}{
		{
			name: "Successfully get all users",
			setupService: func(service *MockUserService) {
				user1 := entities.NewUser("John Doe", "john@example.com")
				user2 := entities.NewUser("Jane Smith", "jane@example.com")
				service.users[user1.ID] = user1
				service.users[user2.ID] = user2
			},
			expectedStatus: http.StatusOK,
			expectedUsers:  2,
		},
		{
			name: "Get users from empty service",
			setupService: func(service *MockUserService) {
				// No users
			},
			expectedStatus: http.StatusOK,
			expectedUsers:  0,
		},
		{
			name: "Service error",
			setupService: func(service *MockUserService) {
				service.shouldFailGet = true
			},
			expectedStatus: http.StatusInternalServerError,
			expectedUsers:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := NewMockUserService()
			tt.setupService(mockService)
			handler := NewUserHandler(mockService)

			req, err := http.NewRequest("GET", "/users", nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			handler.GetUsers(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, status)
			}

			if tt.expectedStatus == http.StatusOK {
				var users []*entities.User
				err := json.NewDecoder(rr.Body).Decode(&users)
				if err != nil {
					t.Errorf("Error decoding response: %v", err)
				}

				if len(users) != tt.expectedUsers {
					t.Errorf("Expected %d users, got %d", tt.expectedUsers, len(users))
				}

				contentType := rr.Header().Get("Content-Type")
				if contentType != "application/json" {
					t.Errorf("Expected Content-Type 'application/json', got '%s'", contentType)
				}
			}
		})
	}
}

func TestUserHandler_GetUser(t *testing.T) {
	tests := []struct {
		name           string
		userID         string
		setupService   func(*MockUserService)
		expectedStatus int
		expectedUserID string
	}{
		{
			name:   "Successfully get existing user",
			userID: "123e4567-e89b-12d3-a456-426614174000",
			setupService: func(service *MockUserService) {
				user := entities.NewUserWithID("123e4567-e89b-12d3-a456-426614174000", "John Doe", "john@example.com")
				service.users["123e4567-e89b-12d3-a456-426614174000"] = user
			},
			expectedStatus: http.StatusOK,
			expectedUserID: "123e4567-e89b-12d3-a456-426614174000",
		},
		{
			name:   "Get non-existing user",
			userID: "nonexistent-id",
			setupService: func(service *MockUserService) {
				// No users
			},
			expectedStatus: http.StatusNotFound,
			expectedUserID: "",
		},
		{
			name:   "Invalid user ID - empty",
			userID: "",
			setupService: func(service *MockUserService) {
				// No setup needed
			},
			expectedStatus: http.StatusBadRequest,
			expectedUserID: "",
		},
		{
			name:   "Service error",
			userID: "123e4567-e89b-12d3-a456-426614174000",
			setupService: func(service *MockUserService) {
				service.shouldFailGet = true
			},
			expectedStatus: http.StatusInternalServerError,
			expectedUserID: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := NewMockUserService()
			tt.setupService(mockService)
			handler := NewUserHandler(mockService)

			req, err := http.NewRequest("GET", "/users/"+tt.userID, nil)
			if err != nil {
				t.Fatal(err)
			}

			// Set up mux vars
			req = mux.SetURLVars(req, map[string]string{"id": tt.userID})

			rr := httptest.NewRecorder()
			handler.GetUser(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, status)
			}

			if tt.expectedStatus == http.StatusOK {
				var user entities.User
				err := json.NewDecoder(rr.Body).Decode(&user)
				if err != nil {
					t.Errorf("Error decoding response: %v", err)
				}

				if user.ID != tt.expectedUserID {
					t.Errorf("Expected user ID %s, got %s", tt.expectedUserID, user.ID)
				}

				contentType := rr.Header().Get("Content-Type")
				if contentType != "application/json" {
					t.Errorf("Expected Content-Type 'application/json', got '%s'", contentType)
				}
			}
		})
	}
}

func TestUserHandler_CreateUser(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    interface{}
		setupService   func(*MockUserService)
		expectedStatus int
		expectedName   string
		expectedEmail  string
	}{
		{
			name: "Successfully create user",
			requestBody: map[string]string{
				"name":  "John Doe",
				"email": "john@example.com",
			},
			setupService:   func(service *MockUserService) {},
			expectedStatus: http.StatusCreated,
			expectedName:   "John Doe",
			expectedEmail:  "john@example.com",
		},
		{
			name: "Create user with empty name",
			requestBody: map[string]string{
				"name":  "",
				"email": "john@example.com",
			},
			setupService:   func(service *MockUserService) {},
			expectedStatus: http.StatusBadRequest,
			expectedName:   "",
			expectedEmail:  "",
		},
		{
			name: "Create user with empty email",
			requestBody: map[string]string{
				"name":  "John Doe",
				"email": "",
			},
			setupService:   func(service *MockUserService) {},
			expectedStatus: http.StatusBadRequest,
			expectedName:   "",
			expectedEmail:  "",
		},
		{
			name:           "Invalid JSON",
			requestBody:    "invalid json",
			setupService:   func(service *MockUserService) {},
			expectedStatus: http.StatusBadRequest,
			expectedName:   "",
			expectedEmail:  "",
		},
		{
			name: "Service error",
			requestBody: map[string]string{
				"name":  "John Doe",
				"email": "john@example.com",
			},
			setupService: func(service *MockUserService) {
				service.shouldFailCreate = true
			},
			expectedStatus: http.StatusInternalServerError,
			expectedName:   "",
			expectedEmail:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := NewMockUserService()
			tt.setupService(mockService)
			handler := NewUserHandler(mockService)

			var body bytes.Buffer
			if str, ok := tt.requestBody.(string); ok {
				body.WriteString(str)
			} else {
				json.NewEncoder(&body).Encode(tt.requestBody)
			}

			req, err := http.NewRequest("POST", "/users", &body)
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			handler.CreateUser(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, status)
			}

			if tt.expectedStatus == http.StatusCreated {
				var user entities.User
				err := json.NewDecoder(rr.Body).Decode(&user)
				if err != nil {
					t.Errorf("Error decoding response: %v", err)
				}

				if user.Name != tt.expectedName {
					t.Errorf("Expected user name '%s', got '%s'", tt.expectedName, user.Name)
				}
				if user.Email != tt.expectedEmail {
					t.Errorf("Expected user email '%s', got '%s'", tt.expectedEmail, user.Email)
				}
				if user.ID == "" {
					t.Error("Expected non-empty user ID")
				}

				contentType := rr.Header().Get("Content-Type")
				if contentType != "application/json" {
					t.Errorf("Expected Content-Type 'application/json', got '%s'", contentType)
				}
			}
		})
	}
}

func TestUserHandler_UpdateUser(t *testing.T) {
	tests := []struct {
		name           string
		userID         string
		requestBody    interface{}
		setupService   func(*MockUserService)
		expectedStatus int
		expectedName   string
		expectedEmail  string
	}{
		{
			name:   "Successfully update user",
			userID: "123e4567-e89b-12d3-a456-426614174000",
			requestBody: map[string]string{
				"name":  "Updated Name",
				"email": "updated@example.com",
			},
			setupService: func(service *MockUserService) {
				user := entities.NewUserWithID("123e4567-e89b-12d3-a456-426614174000", "Original Name", "original@example.com")
				service.users["123e4567-e89b-12d3-a456-426614174000"] = user
			},
			expectedStatus: http.StatusOK,
			expectedName:   "Updated Name",
			expectedEmail:  "updated@example.com",
		},
		{
			name:   "Update non-existing user",
			userID: "nonexistent-id",
			requestBody: map[string]string{
				"name":  "Updated Name",
				"email": "updated@example.com",
			},
			setupService:   func(service *MockUserService) {},
			expectedStatus: http.StatusNotFound,
			expectedName:   "",
			expectedEmail:  "",
		},
		{
			name:   "Invalid user ID - empty",
			userID: "",
			requestBody: map[string]string{
				"name":  "Updated Name",
				"email": "updated@example.com",
			},
			setupService:   func(service *MockUserService) {},
			expectedStatus: http.StatusBadRequest,
			expectedName:   "",
			expectedEmail:  "",
		},
		{
			name:   "Update user with empty name",
			userID: "123e4567-e89b-12d3-a456-426614174000",
			requestBody: map[string]string{
				"name":  "",
				"email": "updated@example.com",
			},
			setupService: func(service *MockUserService) {
				user := entities.NewUserWithID("123e4567-e89b-12d3-a456-426614174000", "Original Name", "original@example.com")
				service.users["123e4567-e89b-12d3-a456-426614174000"] = user
			},
			expectedStatus: http.StatusBadRequest,
			expectedName:   "",
			expectedEmail:  "",
		},
		{
			name:   "Update user with empty email",
			userID: "123e4567-e89b-12d3-a456-426614174000",
			requestBody: map[string]string{
				"name":  "Updated Name",
				"email": "",
			},
			setupService: func(service *MockUserService) {
				user := entities.NewUserWithID("123e4567-e89b-12d3-a456-426614174000", "Original Name", "original@example.com")
				service.users["123e4567-e89b-12d3-a456-426614174000"] = user
			},
			expectedStatus: http.StatusBadRequest,
			expectedName:   "",
			expectedEmail:  "",
		},
		{
			name:           "Invalid JSON",
			userID:         "123e4567-e89b-12d3-a456-426614174000",
			requestBody:    "invalid json",
			setupService:   func(service *MockUserService) {},
			expectedStatus: http.StatusBadRequest,
			expectedName:   "",
			expectedEmail:  "",
		},
		{
			name:   "Service error",
			userID: "123e4567-e89b-12d3-a456-426614174000",
			requestBody: map[string]string{
				"name":  "Updated Name",
				"email": "updated@example.com",
			},
			setupService: func(service *MockUserService) {
				service.shouldFailUpdate = true
			},
			expectedStatus: http.StatusInternalServerError,
			expectedName:   "",
			expectedEmail:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := NewMockUserService()
			tt.setupService(mockService)
			handler := NewUserHandler(mockService)

			var body bytes.Buffer
			if str, ok := tt.requestBody.(string); ok {
				body.WriteString(str)
			} else {
				json.NewEncoder(&body).Encode(tt.requestBody)
			}

			req, err := http.NewRequest("PUT", "/users/"+tt.userID, &body)
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Content-Type", "application/json")

			// Set up mux vars
			req = mux.SetURLVars(req, map[string]string{"id": tt.userID})

			rr := httptest.NewRecorder()
			handler.UpdateUser(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, status)
			}

			if tt.expectedStatus == http.StatusOK {
				var user entities.User
				err := json.NewDecoder(rr.Body).Decode(&user)
				if err != nil {
					t.Errorf("Error decoding response: %v", err)
				}

				if user.Name != tt.expectedName {
					t.Errorf("Expected user name '%s', got '%s'", tt.expectedName, user.Name)
				}
				if user.Email != tt.expectedEmail {
					t.Errorf("Expected user email '%s', got '%s'", tt.expectedEmail, user.Email)
				}

				contentType := rr.Header().Get("Content-Type")
				if contentType != "application/json" {
					t.Errorf("Expected Content-Type 'application/json', got '%s'", contentType)
				}
			}
		})
	}
}

func TestUserHandler_DeleteUser(t *testing.T) {
	tests := []struct {
		name           string
		userID         string
		setupService   func(*MockUserService)
		expectedStatus int
	}{
		{
			name:   "Successfully delete user",
			userID: "123e4567-e89b-12d3-a456-426614174000",
			setupService: func(service *MockUserService) {
				user := entities.NewUserWithID("123e4567-e89b-12d3-a456-426614174000", "John Doe", "john@example.com")
				service.users["123e4567-e89b-12d3-a456-426614174000"] = user
			},
			expectedStatus: http.StatusNoContent,
		},
		{
			name:           "Delete non-existing user",
			userID:         "nonexistent-id",
			setupService:   func(service *MockUserService) {},
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "Invalid user ID - empty",
			userID:         "",
			setupService:   func(service *MockUserService) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "Service error",
			userID: "123e4567-e89b-12d3-a456-426614174000",
			setupService: func(service *MockUserService) {
				service.shouldFailDelete = true
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := NewMockUserService()
			tt.setupService(mockService)
			handler := NewUserHandler(mockService)

			req, err := http.NewRequest("DELETE", "/users/"+tt.userID, nil)
			if err != nil {
				t.Fatal(err)
			}

			// Set up mux vars
			req = mux.SetURLVars(req, map[string]string{"id": tt.userID})

			rr := httptest.NewRecorder()
			handler.DeleteUser(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, status)
			}

			// For successful deletion, check that response body is empty
			if tt.expectedStatus == http.StatusNoContent {
				responseBody := strings.TrimSpace(rr.Body.String())
				if responseBody != "" {
					t.Errorf("Expected empty response body for successful deletion, got '%s'", responseBody)
				}
			}
		})
	}
}
