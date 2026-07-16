CREATE TABLE items (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name VARCHAR(255) UNIQUE NOT NULL,
  description VARCHAR(255) NOT NULL,
  category_id INTEGER NOT NULL REFERENCES item_categories(id),
  effect_target VARCHAR(255),
  effect_value INTEGER,
  value INTEGER NOT NULL
);