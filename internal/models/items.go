package models

import "errors"

type Effect struct {
	Description string
	Target      string
	Value       int32
}

type Item struct {
	name        string
	description string
	category    string
	effect      Effect
	value       int32
}

func NewItem(name, description, category string, effect Effect, value int32) (*Item, error) {
	if name == "" {
		return nil, errors.New("invalid item name")
	}
	if description == "" {
		return nil, errors.New("invalid item description")
	}
	if category == "" {
		return nil, errors.New("invalid item category")
	}
	if value < 0 {
		return nil, errors.New("invalid item value")
	}

	item := Item{
		name:        name,
		description: description,
		category:    category,
		effect:      effect,
		value:       value,
	}

	return &item, nil
}

func (i *Item) GetName() string {
	return i.name
}

func (i *Item) GetDescription() string {
	return i.description
}

func (i *Item) GetCategory() string {
	return i.category
}

func (i *Item) GetEffect() Effect {
	return i.effect
}

func (i *Item) GetValue() int32 {
	return i.value
}
