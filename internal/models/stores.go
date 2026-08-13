package models

import (
	"errors"
	"fmt"
	"math/rand/v2"
	"time"
)

type Store struct {
	locationID   int32
	name         string
	gold         int32
	inventory    map[int32]int32 // map[itemID]quantity
	last_updated time.Time
}

func InitStore(locationID int32, name string, itemList map[int32]*Item) (*Store, error) {
	if locationID < 1 {
		return nil, errors.New("invalid location id")
	}
	if name == "" {
		return nil, errors.New("invalid store name")
	}
	store := Store{
		locationID:   locationID,
		name:         name,
		gold:         1000,
		inventory:    generateInventory(itemList),
		last_updated: time.Now().UTC(),
	}
	return &store, nil
}

func generateInventory(items map[int32]*Item) map[int32]int32 {
	inventory := make(map[int32]int32)
	for key := range items {
		// #nosec G404 -- false positive: random number is for generating inventory. Not critical.
		randNum := rand.IntN(100) + 1
		if randNum > 49 {
			// #nosec G404 -- false positive: random number is for generating inventory. Not critical.
			randAmount := rand.IntN(10) + 1
			// #nosec G115 -- false positive: randAmount never going to be higher than 10
			inventory[key] = int32(randAmount)
		}
	}
	return inventory
}

// Basic Functions
func (s *Store) GetLocationID() int32 {
	return s.locationID
}

func (s *Store) GetName() string {
	return s.name
}

// Gold Functions
func (s *Store) GetGold() int32 {
	return s.gold
}

func (s *Store) IncreaseGold(amount int32) error {
	if amount < 1 {
		return errors.New("invalid gold amount")
	}
	s.gold += amount
	s.last_updated = time.Now().UTC()
	return nil
}

func (s *Store) DecreaseGold(amount int32) error {
	if amount < 1 || amount > s.gold {
		return errors.New("invalid gold amount")
	}
	s.gold -= amount
	s.last_updated = time.Now().UTC()
	return nil
}

// Inventory Functions
func (s *Store) GetInventory() map[int32]int32 {
	return s.inventory
}

func (s *Store) RefreshStore(itemList map[int32]*Item) {
	refreshTime := s.last_updated.Add(time.Hour)
	if time.Now().UTC().After(refreshTime) {
		s.gold = 1000
		s.inventory = generateInventory(itemList)
		s.last_updated = time.Now().UTC()
	}
}

func (s *Store) BuyItem(itemList map[int32]*Item, itemID, quantity int32, player *Player) (n int32, err error) {
	// Check Item Exists
	item, ok := itemList[itemID]
	if itemID < 1 || !ok {
		return 0, errors.New("invalid item id")
	}
	if quantity < 1 {
		return 0, errors.New("invalid quantity")
	}

	// Check Store Inventory
	storeQuantity, ok := s.inventory[itemID]
	if !ok {
		//lint:ignore ST1005 NPC talking
		return 0, fmt.Errorf("I'm sorry. I don't carry %s.", itemList[itemID].GetName()) //lint:ignore ST1005 NPC talking
	}
	if quantity > storeQuantity {
		//lint:ignore ST1005 NPC talking
		return 0, fmt.Errorf("I'm sorry. I only have %d in stock.", storeQuantity)
	}

	// Get total amount
	totalAmount := item.GetValue() * quantity

	// Gold Transfer
	err = player.RemoveGold(totalAmount)
	if err != nil {
		return 0, err
	}
	err = s.IncreaseGold(totalAmount)
	if err != nil {
		return 0, err
	}

	// Item Transfer
	err = player.AddItem(itemID, quantity)
	if err != nil {
		return 0, err
	}
	if storeQuantity == quantity {
		delete(s.inventory, itemID)
	} else {
		s.inventory[itemID] -= quantity
	}

	// Return
	return quantity, nil
}

func (s *Store) SellItem(itemList map[int32]*Item, itemID, quantity int32, player *Player) (int32, error) {
	// Check Item Exist
	item, ok := itemList[itemID]
	if itemID < 1 || !ok {
		return 0, errors.New("invalid item ID")
	}
	if quantity < 1 {
		return 0, errors.New("invalid quantity")
	}

	// Check Player Inventory
	amount, ok := player.GetInventory()[itemID]
	if !ok {
		//lint:ignore ST1005 NPC talking
		return 0, fmt.Errorf("Oh. Looks like you don't have any %s.", item.GetName())
	}
	if amount < quantity {
		//lint:ignore ST1005 NPC talking
		return 0, fmt.Errorf("Oh. Looks like you don't only have %d in your inventory.", amount)
	}

	// Get Total Amount
	totalAmount := item.GetValue() * quantity

	// Gold Transfer
	if totalAmount > s.GetGold() {
		// Store always wins
		totalAmount = s.GetGold()
	}
	err := s.DecreaseGold(totalAmount)
	if err != nil {
		return 0, err
	}
	err = player.AddGold(totalAmount)
	if err != nil {
		return 0, err
	}

	// Item Transfer
	_, ok = s.inventory[itemID]
	if !ok {
		s.inventory[itemID] = quantity
	} else {
		s.inventory[itemID] += quantity
	}
	_, err = player.RemoveItem(itemID, quantity)
	if err != nil {
		return 0, err
	}

	// Return
	return totalAmount, nil
}
