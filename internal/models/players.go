package models

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Player struct {
	id        uuid.UUID
	name      string
	level     level
	health    health
	mana      mana
	stat      stat
	buff      *Buff
	inventory map[int32]int32 // map[itemID]quantity
	gold      int32
	location  int32
	equipment equipment
}

type level struct {
	currentExp   int32
	currentLevel int32
}

type Buff struct {
	target   string
	value    int32
	expireAt time.Time
}

type equipment struct {
	weapon    int32
	armor     int32
	accessory int32
}

// Initialization
func InitPlayer(id uuid.UUID, name string, currentExp, currentLevel, gold, location int32, inventory map[int32]int32, weapon, armor, accessory int32) (*Player, error) {
	var uuidNil uuid.UUID
	if id == uuidNil {
		return nil, errors.New("invalid id")
	}
	if name == "" {
		return nil, errors.New("invalid player name")
	}
	player := Player{
		id:   id,
		name: name,
	}

	// Set Level
	if currentExp >= 0 || currentExp <= 1000 {
		player.level.currentExp = currentExp
	}
	player.level.currentLevel = max(currentLevel, 1)

	// Set Health
	player.health.maxHealth = ((player.level.currentLevel - 1) * 10) + 100
	player.health.currentHealth = player.health.maxHealth

	// Set Mana
	player.mana.maxMana = ((player.level.currentLevel - 1) * 10) + 100
	player.mana.currentMana = player.mana.maxMana

	// Set Stat
	player.stat.strength = ((player.level.currentLevel - 1) * 2) + 10
	player.stat.defense = ((player.level.currentLevel - 1) * 2) + 10

	// Set Inventory
	player.inventory = inventory

	// Set Gold
	player.gold = max(gold, 0)

	// Set Location
	player.location = location

	// Set Equipment
	player.equipment = equipment{
		weapon:    weapon,
		armor:     armor,
		accessory: accessory,
	}

	return &player, nil
}

// Name Function
func (p *Player) GetName() string {
	return p.name
}

// Level Functions
func (p *Player) GetLevel() int32 {
	return p.level.currentLevel
}

func (p *Player) GetCurrentExp() int32 {
	return p.level.currentExp
}

func (p *Player) AddExp(amount int32) (bool, error) {
	if amount < 0 {
		return false, errors.New("invalid exp amount")
	}
	p.level.currentExp += amount
	if p.level.currentExp >= 1000 {

		// Increase Level
		p.level.currentLevel += 1
		p.level.currentExp -= 1000

		// Increase Health
		p.health.maxHealth += 10
		p.health.currentHealth = p.health.maxHealth

		// Increase Mana
		p.mana.maxMana += 10
		p.mana.currentMana = p.mana.maxMana

		// Increase Stats
		p.stat.strength += 2
		p.stat.defense += 2

		return true, nil
	}
	return false, nil
}

// Health Functions
func (p *Player) GetCurrentHealth() int32 {
	return p.health.currentHealth
}

func (p *Player) GetMaxHealth() int32 {
	return p.health.maxHealth
}

func (p *Player) TakeDamage(damage int32) (bool, error) {
	if damage < 0 {
		return false, errors.New("invalid damage")
	}
	p.health.currentHealth -= damage
	if p.health.currentHealth <= 0 {
		p.health.currentHealth = 0

		// Lose Exp
		p.level.currentExp -= 50
		if p.level.currentExp < 0 {
			p.level.currentExp = 0
		}
		return true, nil
	}
	return false, nil
}

func (p *Player) RestoreHealth(amount int32) error {
	if amount < 0 {
		return errors.New("invalid health amount")
	}
	p.health.currentHealth += amount
	if p.health.currentHealth > p.health.maxHealth {
		p.health.currentHealth = p.health.maxHealth
	}
	return nil
}

// Mana Functions
func (p *Player) GetCurrentMana() int32 {
	return p.mana.currentMana
}

func (p *Player) GetMaxMana() int32 {
	return p.mana.maxMana
}

func (p *Player) RestoreMana(amount int32) error {
	if amount < 0 {
		return errors.New("invalid mana amount")
	}
	p.mana.currentMana += amount
	if p.mana.currentMana > p.mana.maxMana {
		p.mana.currentMana = p.mana.maxMana
	}
	return nil
}

func (p *Player) UseMana(amount int32) error {
	if amount > p.mana.currentMana {
		return errors.New("not enough mana")
	}
	p.mana.currentMana -= amount
	return nil
}

// Stat Functions
func (p *Player) GetStrength() int32 {
	var buffValue int32 = 0
	if p.buff != nil && p.buff.target == "Strength" {
		buffValue = p.buff.value
	}
	return p.stat.strength + buffValue
}

func (p *Player) GetDefense() int32 {
	var buffValue int32 = 0
	if p.buff != nil && p.buff.target == "Defense" {
		buffValue = p.buff.value
	}
	return p.stat.defense + buffValue
}

// Buff Functions
func (p *Player) GetBuff() (*Buff, error) {
	if p.buff == nil || time.Now().UTC().After(p.buff.expireAt) {
		switch p.buff.target {
		case "Health":
			p.health.maxHealth -= p.buff.value
			if p.health.currentHealth > p.health.maxHealth {
				p.health.currentHealth = p.health.maxHealth
			}
		case "Mana":
			p.mana.maxMana -= p.buff.value
			if p.mana.currentMana > p.mana.maxMana {
				p.mana.currentMana = p.mana.maxMana
			}
		}
		p.buff = nil
		return nil, errors.New("no active buff")
	}
	return p.buff, nil
}

func (p *Player) EatFood(itemList map[int32]*Item, itemID int32) error {
	if itemID < 1 {
		return errors.New("invalid item ID")
	}
	if itemList[itemID].category != "Meal" {
		return errors.New("item is not a food item")
	}
	quantity, ok := p.inventory[itemID]
	if !ok {
		return errors.New("you don't have that item")
	}
	item := itemList[itemID]
	buff := Buff{
		target:   item.GetEffect().Target,
		value:    item.GetEffect().Value,
		expireAt: time.Now().UTC().Add(30 * time.Minute),
	}
	p.buff = &buff
	if quantity == 1 {
		delete(p.inventory, itemID)
	} else {
		p.inventory[itemID] -= 1
	}
	switch buff.target {
	case "Health":
		p.health.maxHealth += buff.value
		p.health.currentHealth += buff.value
	case "Mana":
		p.mana.maxMana += buff.value
		p.mana.currentMana += buff.value
	}
	fmt.Printf("You ate %s.\n", item.GetName())
	fmt.Printf("Your %s was increase by %d for 30 minutes.\n", buff.target, buff.value)
	return nil
}

// Inventory Functions
func (p *Player) GetInventory() map[int32]int32 {
	return p.inventory
}

func (p *Player) AddItem(itemID int32, amount int32) error {
	if amount <= 0 {
		return errors.New("invalid item amount")
	}
	_, ok := p.inventory[itemID]
	if !ok {
		p.inventory[itemID] = amount
	} else {
		p.inventory[itemID] += amount
	}
	return nil
}

func (p *Player) RemoveItem(itemID int32, amount int32) (int32, error) {
	num, ok := p.inventory[itemID]
	if !ok {
		return 0, errors.New("invalid item")
	}
	if num <= amount {
		delete(p.inventory, itemID)
		return num, nil
	}
	p.inventory[itemID] -= amount
	return amount, nil
}

// Gold Functions
func (p *Player) GetGold() int32 {
	return p.gold
}

func (p *Player) AddGold(amount int32) error {
	if amount < 0 {
		return errors.New("invalid gold amount")
	}
	p.gold += amount
	return nil
}

func (p *Player) RemoveGold(amount int32) error {
	if amount > p.gold {
		return errors.New("not enough gold")
	}
	p.gold -= amount
	return nil
}

// Location Functions
func (p *Player) GetLocation() int32 {
	return p.location
}

func (p *Player) SetLocation(location int32) error {
	if location < 1 {
		return errors.New("invalid location")
	}
	p.location = location
	return nil
}

// Display Stats
func (p *Player) DisplayStats() {
	fmt.Println("\n==== Player Stat ====")
	fmt.Printf(" Level: %d\n", p.GetLevel())
	fmt.Printf(" Exp: %d/1000\n", p.GetCurrentExp())
	fmt.Printf(" Health: %d/%d\n", p.GetCurrentHealth(), p.GetMaxHealth())
	fmt.Printf(" Mana: %d/%d\n", p.GetCurrentMana(), p.GetMaxMana())
	fmt.Printf(" Strength: %d\n", p.GetStrength())
	fmt.Printf(" Defense: %d\n", p.GetDefense())
	fmt.Printf(" Gold: %d\n", p.GetGold())
	fmt.Println()
}

func (p *Player) DisplayGear(itemList map[int32]*Item) {
	fmt.Println("\n==== Player Gear ====")
	if p.GetWeapon() == 0 {
		fmt.Println(" Weapon: Not equipped")
	} else {
		fmt.Printf(" Weapon: %s\n", itemList[p.GetWeapon()].GetName())
	}
	if p.GetArmor() == 0 {
		fmt.Println(" Armor: Not equipped")
	} else {
		fmt.Printf(" Armor: %s\n", itemList[p.GetArmor()].GetName())
	}
	if p.GetAccessory() == 0 {
		fmt.Println(" Accessory: Not equipped")
	} else {
		fmt.Printf(" Accessory: %s\n", itemList[p.GetAccessory()].GetName())
	}
	fmt.Println()
}

// Equipment Functions
func (p *Player) GetWeapon() int32 {
	return p.equipment.weapon
}

func (p *Player) GetArmor() int32 {
	return p.equipment.armor
}

func (p *Player) GetAccessory() int32 {
	return p.equipment.accessory
}

func (p *Player) EquipWeapon(itemList map[int32]*Item, weaponID int32) error {
	// Check weapon ID is valid
	if weaponID < 1 {
		return errors.New("invalid weapon ID")
	}
	item, ok := itemList[weaponID]
	if !ok {
		return errors.New("invalid weapon ID")
	}
	if item.GetCategory() != "Weapon" {
		return errors.New("weapon ID is not a weapon")
	}
	// Check if player has weapon in inventory
	_, ok = p.GetInventory()[weaponID]
	if !ok {
		return errors.New("you do not have that weapon")
	}
	// Check if player already has a weapon equipped
	_ = p.UnequipWeapon()

	p.equipment.weapon = weaponID
	_, err := p.RemoveItem(weaponID, 1)
	if err != nil {
		return err
	}
	return nil
}

func (p *Player) UnequipWeapon() error {
	// Check weapon actually equipped
	if p.GetWeapon() == 0 {
		return errors.New("no weapon equipped")
	}
	// Move weapon into inventory
	err := p.AddItem(p.GetWeapon(), 1)
	if err != nil {
		return err
	}
	p.equipment.weapon = 0
	return nil
}

func (p *Player) EquipArmor(itemList map[int32]*Item, armorID int32) error {
	// Check armor ID is valid
	if armorID < 1 {
		return errors.New("invalid armor ID")
	}
	item, ok := itemList[armorID]
	if !ok {
		return errors.New("invalid armor ID")
	}
	if item.GetCategory() != "Armor" {
		return errors.New("armor ID is not an armor")
	}
	// Check if player has armor in inventory
	_, ok = p.GetInventory()[armorID]
	if !ok {
		return errors.New("you do not have that armor")
	}
	// Check if player already has an armor equipped
	_ = p.UnequipArmor()

	p.equipment.armor = armorID
	_, err := p.RemoveItem(armorID, 1)
	if err != nil {
		return err
	}
	return nil
}

func (p *Player) UnequipArmor() error {
	// Check armor actually equipped
	if p.GetArmor() == 0 {
		return errors.New("no armor equipped")
	}
	// Move armor into inventory
	err := p.AddItem(p.GetArmor(), 1)
	if err != nil {
		return err
	}
	p.equipment.armor = 0
	return nil
}

func (p *Player) EquipAccessory(itemList map[int32]*Item, accessoryID int32) error {
	// Check accessory ID is valid
	if accessoryID < 1 {
		return errors.New("invalid accessory ID")
	}
	item, ok := itemList[accessoryID]
	if !ok {
		return errors.New("invalid accessory ID")
	}
	if item.GetCategory() != "Accessory" {
		return errors.New("accessory ID is not an accessory")
	}
	// Check if player has accessory in inventory
	_, ok = p.GetInventory()[accessoryID]
	if !ok {
		return errors.New("you do not have that accessory")
	}
	// Check if player already has an accessory equipped
	_ = p.UnequipAccessory()

	p.equipment.accessory = accessoryID
	_, err := p.RemoveItem(accessoryID, 1)
	if err != nil {
		return err
	}
	return nil
}

func (p *Player) UnequipAccessory() error {
	// Check accessory actually equipped
	if p.GetAccessory() == 0 {
		return errors.New("no accessory equipped")
	}
	// Move accessory into inventory
	err := p.AddItem(p.GetAccessory(), 1)
	if err != nil {
		return err
	}
	p.equipment.accessory = 0
	return nil
}
