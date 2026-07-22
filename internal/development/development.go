package development

import (
	"context"
	"fmt"

	"github.com/poupardm-GhostWrath/GoAdventure/internal/auth"
	"github.com/poupardm-GhostWrath/GoAdventure/internal/database"
)

const (
	testEmail = "test@example.com"
	testPass  = "test"
)

func CreateTestUser(ctx context.Context, dbQueries *database.Queries) error {
	_, err := dbQueries.GetUserByEmail(ctx, testEmail)
	if err != nil {
		hash, err := auth.HashPassword(testPass)
		if err != nil {
			return fmt.Errorf("failed to hash password for test user: %v", err)
		}
		err = dbQueries.CreateUser(ctx, database.CreateUserParams{
			Email:        testEmail,
			PasswordHash: hash,
		})
		if err != nil {
			return fmt.Errorf("failed to create test user in DB: %v", err)
		}
	}
	return nil
}
