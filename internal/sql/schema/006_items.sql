CREATE TABLE items (
  id SERIAL PRIMARY KEY,
  name VARCHAR(255) UNIQUE NOT NULL,
  description VARCHAR(255) NOT NULL,
  category_id INTEGER NOT NULL REFERENCES item_categories(id),
  effect_description VARCHAR(255),
  effect_target VARCHAR(255),
  effect_value INTEGER,
  value INTEGER NOT NULL
);