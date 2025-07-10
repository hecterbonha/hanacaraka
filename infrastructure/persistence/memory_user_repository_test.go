package persistence

import (
	"sync"
	"testing"

	"com.hanacaraka/domain/entities"
)

func TestNewMemoryUserRepository(t *testing.T) {
	repo := NewMemoryUserRepository()

	if repo == nil {
		t.Error("Expected repository to be created, got nil")
	}

	// Test that it returns the interface type
	_, ok := repo.(*MemoryUserRepository)
	if !ok {
		t.Error("Expected repository to be of type *MemoryUserRepository")
	}

	// Test that it has initial data
	users, err := repo.GetAll()
	if err != nil {
		t.Errorf("Expected no error getting all users, got %v", err)
	}

	if len(users) != 3 {
		t.Errorf("Expected 3 initial users, got %d", len(users))
	}

	// Verify that users have valid UUIDs and expected data
	expectedNames := []string{"John Doe", "Jane Smith", "Bob Johnson"}
	expectedEmails := []string{"john@example.com", "jane@example.com", "bob@example.com"}

	for i, user := range users {
		if user.ID == "" {
			t.Errorf("Expected user %d to have a non-empty ID", i)
		}
		if i < len(expectedNames) && user.Name != expectedNames[i] {
			t.Errorf("Expected user name '%s', got '%s'", expectedNames[i], user.Name)
		}
		if i < len(expectedEmails) && user.Email != expectedEmails[i] {
			t.Errorf("Expected user email '%s', got '%s'", expectedEmails[i], user.Email)
		}
	}
}

func TestMemoryUserRepository_GetAll(t *testing.T) {
	repo := &MemoryUserRepository{
		users: make([]*entities.User, 0),
	}

	// Test empty repository
	users, err := repo.GetAll()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if len(users) != 0 {
		t.Errorf("Expected 0 users, got %d", len(users))
	}

	// Add some users
	user1 := entities.NewUser("John Doe", "john@example.com")
	user2 := entities.NewUser("Jane Smith", "jane@example.com")
	repo.users = append(repo.users, user1, user2)

	users, err = repo.GetAll()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if len(users) != 2 {
		t.Errorf("Expected 2 users, got %d", len(users))
	}

	// Verify that a copy is returned (not the original slice)
	users[0] = nil
	if repo.users[0] == nil {
		t.Error("Expected original slice to remain unchanged")
	}
}

func TestMemoryUserRepository_GetByID(t *testing.T) {
	repo := &MemoryUserRepository{
		users: make([]*entities.User, 0),
	}

	user1 := entities.NewUser("John Doe", "john@example.com")
	user2 := entities.NewUser("Jane Smith", "jane@example.com")
	repo.users = append(repo.users, user1, user2)

	tests := []struct {
		name          string
		id            string
		expectedUser  *entities.User
		expectedError string
	}{
		{
			name:          "Get existing user",
			id:            user1.ID,
			expectedUser:  user1,
			expectedError: "",
		},
		{
			name:          "Get another existing user",
			id:            user2.ID,
			expectedUser:  user2,
			expectedError: "",
		},
		{
			name:          "Get non-existing user",
			id:            "nonexistent-id",
			expectedUser:  nil,
			expectedError: "user not found",
		},
		{
			name:          "Get user with empty ID",
			id:            "",
			expectedUser:  nil,
			expectedError: "user not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := repo.GetByID(tt.id)

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

func TestMemoryUserRepository_Create(t *testing.T) {
	repo := &MemoryUserRepository{
		users: make([]*entities.User, 0),
	}

	tests := []struct {
		name          string
		user          *entities.User
		expectedName  string
		expectedEmail string
	}{
		{
			name:          "Create first user",
			user:          entities.NewUser("John Doe", "john@example.com"),
			expectedName:  "John Doe",
			expectedEmail: "john@example.com",
		},
		{
			name:          "Create second user",
			user:          entities.NewUser("Jane Smith", "jane@example.com"),
			expectedName:  "Jane Smith",
			expectedEmail: "jane@example.com",
		},
		{
			name:          "Create third user",
			user:          entities.NewUser("Bob Johnson", "bob@example.com"),
			expectedName:  "Bob Johnson",
			expectedEmail: "bob@example.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			initialCount := len(repo.users)

			createdUser, err := repo.Create(tt.user)

			if err != nil {
				t.Errorf("Expected no error, got %v", err)
			}

			if createdUser == nil {
				t.Error("Expected created user, got nil")
				return
			}

			// Check that ID is set and not empty
			if createdUser.ID == "" {
				t.Error("Expected created user to have non-empty ID")
			}

			// Check other fields
			if createdUser.Name != tt.expectedName {
				t.Errorf("Expected name '%s', got '%s'", tt.expectedName, createdUser.Name)
			}
			if createdUser.Email != tt.expectedEmail {
				t.Errorf("Expected email '%s', got '%s'", tt.expectedEmail, createdUser.Email)
			}

			// Check that user was added to repository
			if len(repo.users) != initialCount+1 {
				t.Errorf("Expected %d users in repository, got %d", initialCount+1, len(repo.users))
			}

			// Check that user can be retrieved
			retrievedUser, err := repo.GetByID(createdUser.ID)
			if err != nil {
				t.Errorf("Expected to retrieve created user, got error: %v", err)
			}
			if retrievedUser.ID != createdUser.ID {
				t.Errorf("Retrieved user ID mismatch: expected %s, got %s", createdUser.ID, retrievedUser.ID)
			}
		})
	}
}

func TestMemoryUserRepository_Update(t *testing.T) {
	repo := &MemoryUserRepository{
		users: make([]*entities.User, 0),
	}

	// Add initial users
	user1 := entities.NewUser("John Doe", "john@example.com")
	user2 := entities.NewUser("Jane Smith", "jane@example.com")
	repo.users = append(repo.users, user1, user2)

	tests := []struct {
		name          string
		user          *entities.User
		expectedError string
		shouldExist   bool
	}{
		{
			name:          "Update existing user",
			user:          entities.NewUserWithID(user1.ID, "John Updated", "john.updated@example.com"),
			expectedError: "",
			shouldExist:   true,
		},
		{
			name:          "Update another existing user",
			user:          entities.NewUserWithID(user2.ID, "Jane Updated", "jane.updated@example.com"),
			expectedError: "",
			shouldExist:   true,
		},
		{
			name:          "Update non-existing user",
			user:          entities.NewUserWithID("nonexistent-id", "Non Existing", "nonexisting@example.com"),
			expectedError: "user not found",
			shouldExist:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			initialCount := len(repo.users)

			updatedUser, err := repo.Update(tt.user)

			if tt.expectedError != "" {
				if err == nil {
					t.Errorf("Expected error '%s', got nil", tt.expectedError)
				} else if err.Error() != tt.expectedError {
					t.Errorf("Expected error '%s', got '%s'", tt.expectedError, err.Error())
				}
				if updatedUser != nil {
					t.Errorf("Expected no user to be returned on error, got %+v", updatedUser)
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}

				if updatedUser == nil {
					t.Error("Expected updated user, got nil")
					return
				}

				// Check that user was updated correctly
				if updatedUser.ID != tt.user.ID {
					t.Errorf("Expected ID %s, got %s", tt.user.ID, updatedUser.ID)
				}
				if updatedUser.Name != tt.user.Name {
					t.Errorf("Expected name '%s', got '%s'", tt.user.Name, updatedUser.Name)
				}
				if updatedUser.Email != tt.user.Email {
					t.Errorf("Expected email '%s', got '%s'", tt.user.Email, updatedUser.Email)
				}

				// Check that user can be retrieved with updated values
				retrievedUser, err := repo.GetByID(tt.user.ID)
				if err != nil {
					t.Errorf("Expected to retrieve updated user, got error: %v", err)
				}
				if retrievedUser.Name != tt.user.Name || retrievedUser.Email != tt.user.Email {
					t.Errorf("Retrieved user values not updated correctly")
				}
			}

			// Check that repository size didn't change
			if len(repo.users) != initialCount {
				t.Errorf("Expected repository size to remain %d, got %d", initialCount, len(repo.users))
			}
		})
	}
}

func TestMemoryUserRepository_Delete(t *testing.T) {
	repo := &MemoryUserRepository{
		users: make([]*entities.User, 0),
	}

	// Add initial users
	user1 := entities.NewUser("John Doe", "john@example.com")
	user2 := entities.NewUser("Jane Smith", "jane@example.com")
	user3 := entities.NewUser("Bob Johnson", "bob@example.com")
	repo.users = append(repo.users, user1, user2, user3)

	tests := []struct {
		name          string
		id            string
		expectedError string
		expectedCount int
	}{
		{
			name:          "Delete existing user",
			id:            user2.ID,
			expectedError: "",
			expectedCount: 2,
		},
		{
			name:          "Delete another existing user",
			id:            user1.ID,
			expectedError: "",
			expectedCount: 1,
		},
		{
			name:          "Delete non-existing user",
			id:            "nonexistent-id",
			expectedError: "user not found",
			expectedCount: 1,
		},
		{
			name:          "Delete last user",
			id:            user3.ID,
			expectedError: "",
			expectedCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.Delete(tt.id)

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

				// Check that user was actually deleted
				_, err := repo.GetByID(tt.id)
				if err == nil {
					t.Errorf("Expected user %s to be deleted, but it still exists", tt.id)
				}
			}

			// Check repository size
			if len(repo.users) != tt.expectedCount {
				t.Errorf("Expected %d users in repository, got %d", tt.expectedCount, len(repo.users))
			}
		})
	}
}

// TestMemoryUserRepository_GetNextID is no longer needed since we use UUIDs

func TestMemoryUserRepository_ConcurrentAccess(t *testing.T) {
	repo := NewMemoryUserRepository().(*MemoryUserRepository)

	// Test concurrent reads and writes
	var wg sync.WaitGroup
	numGoroutines := 10
	numOperationsPerGoroutine := 100

	// Concurrent creates
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func(goroutineID int) {
			defer wg.Done()
			for j := 0; j < numOperationsPerGoroutine; j++ {
				user := entities.NewUser("User", "user@example.com")
				_, err := repo.Create(user)
				if err != nil {
					t.Errorf("Goroutine %d: Error creating user: %v", goroutineID, err)
				}
			}
		}(i)
	}

	// Concurrent reads
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func(goroutineID int) {
			defer wg.Done()
			for j := 0; j < numOperationsPerGoroutine; j++ {
				_, err := repo.GetAll()
				if err != nil {
					t.Errorf("Goroutine %d: Error getting all users: %v", goroutineID, err)
				}
			}
		}(i)
	}

	wg.Wait()

	// Check final state
	users, err := repo.GetAll()
	if err != nil {
		t.Errorf("Error getting final user count: %v", err)
	}

	expectedCount := 3 + (numGoroutines * numOperationsPerGoroutine) // 3 initial users + created users
	if len(users) != expectedCount {
		t.Errorf("Expected %d users after concurrent operations, got %d", expectedCount, len(users))
	}
}

func TestMemoryUserRepository_ThreadSafety(t *testing.T) {
	repo := NewMemoryUserRepository().(*MemoryUserRepository)

	var wg sync.WaitGroup
	numGoroutines := 5

	// Test that concurrent operations don't cause data races or panics
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func(goroutineID int) {
			defer wg.Done()

			// Create
			user := entities.NewUser("Concurrent User", "concurrent@example.com")
			createdUser, err := repo.Create(user)
			if err != nil {
				t.Errorf("Goroutine %d: Create error: %v", goroutineID, err)
				return
			}

			// Read
			_, err = repo.GetByID(createdUser.ID)
			if err != nil {
				t.Errorf("Goroutine %d: GetByID error: %v", goroutineID, err)
			}

			// Update
			createdUser.Name = "Updated Name"
			_, err = repo.Update(createdUser)
			if err != nil {
				t.Errorf("Goroutine %d: Update error: %v", goroutineID, err)
			}

			// GetAll
			_, err = repo.GetAll()
			if err != nil {
				t.Errorf("Goroutine %d: GetAll error: %v", goroutineID, err)
			}

			// Delete
			err = repo.Delete(createdUser.ID)
			if err != nil {
				t.Errorf("Goroutine %d: Delete error: %v", goroutineID, err)
			}
		}(i)
	}

	wg.Wait()
}
