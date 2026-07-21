package initialization

import (
	"context"

	"github.com/poupardm-GhostWrath/GoAdventure/internal/database"
	"github.com/poupardm-GhostWrath/GoAdventure/internal/models"
)

func InitializeLocations(dbQueries *database.Queries) (map[int32]*models.Location, error) {
	locations := make(map[int32]*models.Location)
	dbLocations, err := dbQueries.GetLocations(context.Background())
	if err != nil {
		return locations, err
	}
	for _, dbLocation := range dbLocations {
		var directions []models.LocationDirection
		dbDirections, err := dbQueries.GetLocationDirectionByID(context.Background(), dbLocation.ID)
		if err != nil {
			return locations, err
		}
		for _, dbDirection := range dbDirections {
			direction := models.LocationDirection{
				Direction:        dbDirection.Direction,
				TargetLocationID: dbDirection.DirectionTarget,
			}
			directions = append(directions, direction)
		}
		location, err := models.CreateLocation(
			dbLocation.ID,
			dbLocation.Name,
			dbLocation.Description,
			dbLocation.HasStore,
			dbLocation.CanTeleport,
			directions)
		if err != nil {
			return locations, err
		}
		locations[dbLocation.ID] = location
	}
	return locations, nil
}
