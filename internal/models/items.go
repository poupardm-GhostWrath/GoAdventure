package models

type ItemCategory int

const (
	Weapon ItemCategory = iota
	Armor
	Accessory
	Medicine
	Meal
	Ingredient
	Material
	Junk
)

type Effect struct {
	Target string
	Value  int32
}

type Item struct {
	Name        string
	Description string
	Type        ItemCategory
	Effect      Effect
	Value       int32
}
