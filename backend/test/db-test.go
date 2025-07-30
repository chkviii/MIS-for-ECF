package main

import (
	"fmt"
	"mypage-backend/internal/config"
	"mypage-backend/internal/repo"
)

func main() {
	// Load configuration
	cfg := config.Load()
	fmt.Println("GO: Loaded configuration:", *cfg)

	// Initialize database connection
	db, err := repo.InitDB(cfg)
	if err != nil {
		fmt.Println("GO: Error initializing database:", err)
		return
	}

	// Create a new user repository
	userRepo := repo.NewUserRepository(db)

	// Example usage of user repository
	newUser := &repo.User{
		Username: "testuser",
		Email:    "test@test.com",
		Password: "securepassword",
		Role:     "user",
		Status:   "active",
	}
	err = userRepo.Create(newUser)

	if err != nil {
		fmt.Println("GO: Error creating user:", err)
	} else {
		fmt.Println("GO: User created successfully:", newUser)
	}

	newUser2 := &repo.User{
		Username: "testuser2",
	}
	err = userRepo.Create(newUser2)
	if err != nil {
		fmt.Println("GO: Error creating second user:", err)
	} else {
		fmt.Println("GO: Second user created successfully:", newUser2)
	}

	// Example of getting a user by username
	user, err := userRepo.GetByUsername("TeSTuser")
	if err != nil {
		fmt.Println("GO: Error getting user by username:", err)
	} else {
		fmt.Println("GO: Retrieved user:", user)
	}

	// Example of updating a user
	newUser.Username = "testusr3"
	User := newUser
	err = userRepo.Update(1, User)
	if err != nil {
		fmt.Println("GO: Error updating user:", err)
	} else {
		fmt.Println("GO: User updated successfully:", User)
	}

	// Example of find all users
	users, err := userRepo.FindAll()
	if err != nil {
		fmt.Println("GO: Error getting all users:", err)
	} else {
		fmt.Println("GO: Retrieved all users:", users, users[1].Username, users[0].ID)
	}

	//close the database connection
	err = repo.CloseDB(db)
	if err != nil {
		fmt.Println("GO: Error closing database connection:", err)
	} else {
		fmt.Println("GO: Database connection closed successfully")
	}

	//reconnect to the database
	db, err = repo.InitDB(cfg)
	if err != nil {
		fmt.Println("GO: Error re-initializing database:", err)
		return
	}
	fmt.Println("GO: Reconnected to the database successfully")

	// reopen the user repository
	userRepo = repo.NewUserRepository(db)
	// Example of getting a user by ID
	userByID, err := userRepo.GetByID(1)
	if err != nil {
		fmt.Println("GO: Error getting user by ID:", err)
	} else {
		fmt.Println("GO: Retrieved user by ID:", userByID)
	}

}
