CREATE TABLE IF NOT EXISTS recipes (
    id integer PRIMARY KEY,
    name text NOT NULL,
    description text NOT NULL,
    created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    ingredients text NOT NULL default '',
    instructions text NOT NULL default '',
    cooking_time integer NOT NULL default 0,
    image text NOT NULL default ''
);
