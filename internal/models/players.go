package models

import (
	"errors"
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
	inventory map[int32]int32
	gold      int32
	location  int32
}

type level struct {
	currentExp   int32
	currentLevel int32
}

type health struct {
	currentHealth int32
	maxHealth     int32
}

type mana struct {
	currentMana int32
	maxMana     int32
}

type stat struct {
	strength int32
	defense  int32
}

type Buff struct {
	target   string
	value    int32
	expireAt time.Time
}

func NewPlayer(id uuid.UUID, name string, currentExp, currentLevel, gold, location int32, inventory map[int32]int32) (*Player, error) {
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

	return &player, nil
}

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
	return p.stat.strength
}

func (p *Player) GetDefense() int32 {
	return p.stat.defense
}

// Buff Functions
func (p *Player) GetBuff() (*Buff, error) {
	if p.buff == nil || time.Now().UTC().After(p.buff.expireAt) {
		return nil, errors.New("no active buff")
	}
	return p.buff, nil
}

func (p *Player) AddBuff(buff *Buff) error {
	if buff == nil {
		return errors.New("invalid buff")
	}
	p.buff = buff
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
