package auth

import (
	"bufio"
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/howeyc/gopass"
	"github.com/poupardm-GhostWrath/GoAdventure/internal/database"
)

func Login(ctx context.Context, dbQueries *database.Queries, scanner *bufio.Scanner) (uuid.UUID, error) {
	var nullUserID uuid.UUID
	fmt.Println("\n=== Login ===")
	for i := 0; i < 3; i++ {
		// Email
		var email string
		fmt.Print(" Email: ")
		if scanner.Scan() {
			email = scanner.Text()
		}
		if err := scanner.Err(); err != nil {
			return nullUserID, fmt.Errorf("failed to read email: %v", err)
		}

		// Password
		fmt.Print(" Password: ")
		password, err := gopass.GetPasswd()
		if err != nil {
			return nullUserID, fmt.Errorf("failed to read password: %v", err)
		}

		// Get User From DB
		dbUser, err := dbQueries.GetUserByEmail(ctx, email)
		if err != nil {
			fmt.Println("invalid email/password")
			continue
		}

		// Verify Password
		match, err := CheckPasswordHash(string(password), dbUser.PasswordHash)
		if err != nil {
			return nullUserID, fmt.Errorf("failed to check password hash: %v", err)
		}
		if match {
			return dbUser.ID, nil
		}

		fmt.Println("invalid email/password")
	}
	return nullUserID, errors.New("too many attempts")
}
