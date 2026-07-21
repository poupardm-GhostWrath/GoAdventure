CREATE TABLE location_directions (
  location_id INTEGER NOT NULL REFERENCES locations(id),
  direction VARCHAR(255) NOT NULL,
  direction_target INTEGER NOT NULL REFERENCES locations(id),
  PRIMARY KEY(location_id, direction)
);