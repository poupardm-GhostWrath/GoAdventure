package auth

import (
	"bufio"
	"context"
	"errors"
	"fmt"

	"github.com/howeyc/gopass"
	"github.com/poupardm-GhostWrath/GoAdventure/internal/database"
)

func Register(ctx context.Context, dbQueries *database.Queries, scanner *bufio.Scanner) error {
	fmt.Println("\n=== Register ===")
	for i := 0; i < 3; i++ {
		// Email
		var email string
		fmt.Print(" Email: ")
		if scanner.Scan() {
			email = scanner.Text()
		}
		if err := scanner.Err(); err != nil {
			return fmt.Errorf("failed to read email: %v", err)
		}

		// Password
		fmt.Print(" Password: ")
		password, err := gopass.GetPasswd()
		if err != nil {
			return fmt.Errorf("failed to read password: %v", err)
		}

		if email == "" || string(password) == "" {
			fmt.Println("invalid email/password")
			continue
		}

		// Check if email already in use
		_, err = dbQueries.GetUserByEmail(ctx, email)
		if err == nil {
			fmt.Println("email already in use")
			return nil
		}

		// Hash Password
		passwordHash, err := HashPassword(string(password))
		if err != nil {
			return fmt.Errorf("failed to hash password: %v", err)
		}

		err = dbQueries.CreateUser(ctx, database.CreateUserParams{
			Email:        email,
			PasswordHash: passwordHash,
		})
		if err != nil {
			return fmt.Errorf("failed to create user: %v", err)
		} else {
			return nil
		}
	}
	return errors.New("Too many attempts.")
}
