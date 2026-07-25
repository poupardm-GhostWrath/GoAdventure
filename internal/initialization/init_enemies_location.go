package initialization

import (
	"math/rand/v2"

	"github.com/poupardm-GhostWrath/GoAdventure/internal/models"
)

// map[locationID]map[enemyID]quantity
func InitializeEnemiesLocation(enemies map[int32]*models.Enemy, locations map[int32]*models.Location) (map[int32]map[int32]int32, error) {
	var enemiesLocation = make(map[int32]map[int32]int32)
	for _, location := range locations {
		// Check if no store == not town
		if !location.HasStore() {
			var enemiesList = make(map[int32]int32)
			// Go through enemy list
			for _, enemy := range enemies {
				// Check if enemy level is valid for location
				if enemy.GetLevel() >= location.GetMinLevel() && enemy.GetLevel() <= location.GetMaxLevel() {
					// Randomize if enemy swap
					randNum := rand.IntN(100) + 1
					if randNum >= 50 {
						// Randomize enemy quantity
						amountNum := rand.IntN(5) + 1
						enemiesList[enemy.GetID()] = int32(amountNum)
					}
				}
			}
			// Add enemiesList to enemiesLocation
			enemiesLocation[location.GetID()] = enemiesList
		}
	}
	return enemiesLocation, nil
}
