package models

import (
	"errors"
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
		randNum := rand.IntN(100) + 1
		if randNum > 49 {
			randAmount := rand.IntN(10) + 1
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

func (s *Store) RefreshGold() error {
	refreshTime := s.last_updated.Add(time.Hour)
	if time.Now().UTC().After(refreshTime) {
		s.gold = 1000
		return nil
	}
	return errors.New("too early to refresh")
}

// Inventory Functions
func (s *Store) GetInventory() map[int32]int32 {
	return s.inventory
}

func (s *Store) BuyItem(itemList map[int32]*Item, itemID, quantity int32, player *Player) (n int32, err error) {
	var gold int32
	value, ok := s.inventory[itemID]
	if itemID < 1 || !ok {
		return 0, errors.New("invalid item id")
	}
	if quantity < 1 || quantity > value {
		return 0, errors.New("invalid quantity")
	}
	item := itemList[itemID]
	gold = item.GetValue() * quantity
	err = player.RemoveGold(gold)
	if err != nil {
		return 0, err
	}
	if value == quantity {
		delete(s.inventory, itemID)
	} else {
		s.inventory[itemID] -= quantity
	}
	err = s.IncreaseGold(gold)
	if err != nil {
		return 0, err
	}
	return quantity, nil
}

func (s *Store) SellItem(itemList map[int32]*Item, itemID, quantity int32, player *Player) (int32, error) {
	var gold int32
	item, ok := itemList[itemID]
	if itemID < 1 || !ok {
		return 0, errors.New("invalid item ID")
	}
	if quantity < 1 {
		return 0, errors.New("invalid quantity")
	}
	gold = item.GetValue() * quantity
	_, ok = s.inventory[itemID]
	if !ok {
		s.inventory[itemID] = quantity
	} else {
		s.inventory[itemID] += quantity
	}
	if gold > s.GetGold() {
		gold = s.GetGold()
	}
	err := s.DecreaseGold(gold)
	if err != nil {
		return 0, err
	}
	err = player.AddGold(gold)
	if err != nil {
		return 0, err
	}
	return gold, nil
}
