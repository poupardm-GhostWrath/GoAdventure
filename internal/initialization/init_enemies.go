package initialization

import (
	"context"

	"github.com/poupardm-GhostWrath/GoAdventure/internal/database"
	"github.com/poupardm-GhostWrath/GoAdventure/internal/models"
)

func InitializeEnemies(dbQueries *database.Queries) (map[int32]*models.Enemy, error) {
	var enemyList = make(map[int32]*models.Enemy)
	// Get Enemies From DB
	dbEnemies, err := dbQueries.GetEnemies(context.Background())
	if err != nil {
		return enemyList, err
	}
	for _, dbEnemy := range dbEnemies {
		enemy, err := models.InitEnemy(dbEnemy.ID, dbEnemy.Level, dbEnemy.Name)
		if err != nil {
			return enemyList, err
		}
		enemyList[enemy.GetID()] = enemy
	}
	return enemyList, nil
}
