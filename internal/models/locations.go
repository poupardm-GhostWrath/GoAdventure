package models

import "errors"

type Location struct {
	id          int32
	name        string
	description string
	minLevel    int32
	maxLevel    int32
	hasStore    bool
	canTeleport bool
	directions  []*LocationDirection
}

func CreateLocation(id, min_level, max_level int32, name, description string, has_store, can_teleport bool, directions []*LocationDirection) (*Location, error) {
	if id < 1 {
		return nil, errors.New("invalid location id")
	}
	if name == "" {
		return nil, errors.New("invalid location name")
	}
	if description == "" {
		return nil, errors.New("invalid location description")
	}
	if min_level < 1 {
		return nil, errors.New("invalid location min_level")
	}
	if max_level < 1 || max_level < min_level {
		return nil, errors.New("invalid location max_level")
	}
	location := Location{
		id:          id,
		name:        name,
		description: description,
		hasStore:    has_store,
		minLevel:    min_level,
		maxLevel:    max_level,
		canTeleport: can_teleport,
		directions:  directions,
	}
	return &location, nil
}

func (l *Location) GetID() int32 {
	return l.id
}

func (l *Location) GetName() string {
	return l.name
}

func (l *Location) GetDescription() string {
	return l.description
}

func (l *Location) HasStore() bool {
	return l.hasStore
}

func (l *Location) CanTeleport() bool {
	return l.canTeleport
}

func (l *Location) GetDirections() []*LocationDirection {
	return l.directions
}

func (l *Location) GetMinLevel() int32 {
	return l.minLevel
}

func (l *Location) GetMaxLevel() int32 {
	return l.maxLevel
}
