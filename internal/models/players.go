package models

import (
	"errors"
	"time"
)

type Player struct {
	name      string
	level     level
	health    health
	mana      mana
	stat      stat
	buff      *Buff
	inventory map[Item]int32
	gil       int32
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

func NewPlayer(name string) (*Player, error) {
	if name == "" {
		return nil, errors.New("invalid player name")
	}
	player := Player{
		name: name,
		level: level{
			currentExp:   0,
			currentLevel: 1,
		},
		health: health{
			currentHealth: 100,
			maxHealth:     100,
		},
		mana: mana{
			currentMana: 100,
			maxMana:     100,
		},
		stat: stat{
			strength: 10,
			defense:  10,
		},
	}
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
func (p *Player) GetInventory() map[Item]int32 {
	return p.inventory
}

func (p *Player) AddItem(item Item, amount int32) error {
	if amount <= 0 {
		return errors.New("invalid item amount")
	}
	_, ok := p.inventory[item]
	if !ok {
		p.inventory[item] = amount
	} else {
		p.inventory[item] += amount
	}
	return nil
}

func (p *Player) RemoveItem(item Item, amount int32) (int32, error) {
	num, ok := p.inventory[item]
	if !ok {
		return 0, errors.New("invalid item")
	}
	if num <= amount {
		delete(p.inventory, item)
		return num, nil
	}
	p.inventory[item] -= amount
	return amount, nil
}

// Gil Functions
func (p *Player) GetGil() int32 {
	return p.gil
}

func (p *Player) AddGil(amount int32) error {
	if amount < 0 {
		return errors.New("invalid gil amount")
	}
	p.gil += amount
	return nil
}

func (p *Player) RemoveGil(amount int32) error {
	if amount > p.gil {
		return errors.New("not enough gil")
	}
	p.gil -= amount
	return nil
}
