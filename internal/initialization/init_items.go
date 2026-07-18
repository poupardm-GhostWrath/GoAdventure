package initialization

import (
	"context"

	"github.com/poupardm-GhostWrath/GoAdventure/internal/database"
	"github.com/poupardm-GhostWrath/GoAdventure/internal/models"
)

var itemList = make(map[int32]*models.Item)

var ItemCategories = make(map[int32]string)

func InitializeItems(dbQueries *database.Queries) (map[int32]*models.Item, error) {
	// Get Item Categories
	dbItemCategories, err := dbQueries.GetItemCategories(context.Background())
	if err != nil {
		return itemList, err
	}
	for _, dbItemCategory := range dbItemCategories {
		ItemCategories[dbItemCategory.ID] = dbItemCategory.Name
	}

	// Get Items
	dbItems, err := dbQueries.GetItems(context.Background())
	if err != nil {
		return itemList, err
	}
	for _, dbItem := range dbItems {
		effect := models.Effect{}
		if dbItem.EffectDescription.Valid {
			effect.Description = dbItem.EffectDescription.String
			effect.Target = dbItem.EffectTarget.String
			effect.Value = dbItem.EffectValue.Int32
		}
		item, err := models.NewItem(
			dbItem.Name,
			dbItem.Description,
			ItemCategories[dbItem.CategoryID],
			effect,
			dbItem.Value)
		if err != nil {
			return itemList, err
		}

		itemList[dbItem.ID] = item
	}

	return itemList, nil
}
