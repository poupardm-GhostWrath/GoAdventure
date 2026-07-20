package main

import (
	"context"

	"github.com/google/uuid"
)

func getInventory(ctx context.Context, playerID uuid.UUID) (map[int32]int32, error) {
	inventory := make(map[int32]int32)

	dbInventory, err := Cfg.DBQueries.GetInventoryByPlayerID(ctx, playerID)
	if err != nil {
		return inventory, err
	}

	for _, dbItem := range dbInventory {
		inventory[dbItem.ItemID] = dbItem.Quantity
	}

	return inventory, nil
}
