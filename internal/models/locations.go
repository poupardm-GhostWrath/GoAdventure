package models

import "errors"

type Location struct {
	id          int32
	name        string
	description string
	hasStore    bool
	canTeleport bool
	directions  []*LocationDirection
}

func CreateLocation(id int32, name, description string, has_store, can_teleport bool, directions []*LocationDirection) (*Location, error) {
	if id < 1 {
		return nil, errors.New("invalid location id")
	}
	if name == "" {
		return nil, errors.New("invalid location name")
	}
	if description == "" {
		return nil, errors.New("invalid location description")
	}
	location := Location{
		id:          id,
		name:        name,
		description: description,
		hasStore:    has_store,
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
