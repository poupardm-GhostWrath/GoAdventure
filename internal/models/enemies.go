package models

import "errors"

type Enemy struct {
	id     int32
	name   string
	level  int32
	health health
	stat   stat
}

func InitEnemy(id, level int32, name string) (*Enemy, error) {
	if id < 1 {
		return nil, errors.New("invalid enemy id")
	}
	if level < 1 {
		return nil, errors.New("invalid enemy level")
	}
	if name == "" {
		return nil, errors.New("invalid enemy name")
	}
	healthAmount := ((level - 1) * 10) + 50
	statAmount := ((level - 1) * 2) + 5
	enemy := Enemy{
		id:    id,
		name:  name,
		level: level,
		health: health{
			currentHealth: healthAmount,
			maxHealth:     healthAmount,
		},
		stat: stat{
			strength: statAmount,
			defense:  statAmount,
		},
	}
	return &enemy, nil
}

// Get Functions
func (e *Enemy) GetID() int32 {
	return e.id
}

func (e *Enemy) GetName() string {
	return e.name
}

func (e *Enemy) GetLevel() int32 {
	return e.level
}

func (e *Enemy) GetStrength() int32 {
	return e.stat.strength
}

func (e *Enemy) GetDefense() int32 {
	return e.stat.defense
}

func (e *Enemy) GetCurrentHealth() int32 {
	return e.health.currentHealth
}

func (e *Enemy) GetMaxHealth() int32 {
	return e.health.maxHealth
}

// Action Function
func (e *Enemy) TakeDamage(amount int32) (bool, error) {
	if amount < 1 {
		return false, errors.New("invalid damage amount")
	}
	e.health.currentHealth -= amount
	if e.health.currentHealth <= 0 {
		return true, nil
	}
	return false, nil
}
