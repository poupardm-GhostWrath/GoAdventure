CREATE TABLE inventory (
  item_id INTEGER NOT NULL REFERENCES items(id) ON DELETE CASCADE,
  player_id UUID NOT NULL REFERENCES players(id) ON DELETE CASCADE,
  quantity INTEGER NOT NULL DEFAULT 1,
  PRIMARY KEY (item_id, player_id)
);