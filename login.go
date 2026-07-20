package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/howeyc/gopass"
	"github.com/poupardm-GhostWrath/GoAdventure/internal/auth"
	"github.com/poupardm-GhostWrath/GoAdventure/internal/database"
)

func login(ctx context.Context, scanner *bufio.Scanner) (uuid.UUID, error) {
	var userID uuid.UUID
	fmt.Println("==== Login ====")
	for i := 0; i < 3; i++ {
		// Get Email
		var email string
		fmt.Print("Email: ")
		if scanner.Scan() {
			email = scanner.Text()
		}
		if err := scanner.Err(); err != nil {
			return userID, err
		}

		// Get Password
		fmt.Print("Password: ")
		password, err := gopass.GetPasswd()
		if err != nil {
			return userID, err
		}

		// Get User From DB
		dbUser, err := Cfg.DBQueries.GetUserByEmail(ctx, email)
		if err != nil {
			fmt.Println("invalid email/password")
			continue
		}

		// Check Password
		match, err := auth.CheckPasswordHash(string(password), dbUser.PasswordHash)
		if err != nil {
			return userID, err
		}
		if match {
			return dbUser.ID, nil
		}

		fmt.Println("invalid email/password")
	}
	return userID, errors.New("Too many attempts.")
}

func createTestUser(ctx context.Context) error {
	_, err := Cfg.DBQueries.GetUserByEmail(ctx, "test@example.com")
	if err != nil {
		hash, err := auth.HashPassword("test")
		if err != nil {
			return err
		}
		err = Cfg.DBQueries.CreateUser(ctx, database.CreateUserParams{
			Email:        "test@example.com",
			PasswordHash: hash,
		})
		if err != nil {
			return err
		}
	}
	return nil
}
