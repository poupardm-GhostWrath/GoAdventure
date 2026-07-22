package initialization

import (
	"strings"

	"github.com/poupardm-GhostWrath/GoAdventure/internal/models"
)

func InitializeStore(items map[int32]*models.Item, locations map[int32]*models.Location) (map[int32]*models.Store, error) {
	var storeList = make(map[int32]*models.Store)
	for key, location := range locations {
		if location.HasStore() {
			store, err := models.InitStore(key, strings.Join([]string{location.GetName(), "store"}, " "), items)
			if err != nil {
				return storeList, err
			}
			storeList[key] = store
		}
	}
	return storeList, nil
}
