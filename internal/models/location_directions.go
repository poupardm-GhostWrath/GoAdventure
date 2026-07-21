package models

import "errors"

type LocationDirection struct {
	direction  string
	locationID int32
}

func CreateLocationDirection(direction string, locationID int32) (*LocationDirection, error) {
	if direction == "" {
		return nil, errors.New("invalid direction")
	}
	if locationID < 1 {
		return nil, errors.New("invalid location ID")
	}
	return &LocationDirection{direction: direction, locationID: locationID}, nil
}

func (ld *LocationDirection) GetDirection() string {
	return ld.direction
}

func (ld *LocationDirection) GetLocationID() int32 {
	return ld.locationID
}
