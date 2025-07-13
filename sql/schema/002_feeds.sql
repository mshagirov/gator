-- +goose Up
CREATE TABLE feeds (
  id uuid PRIMARY KEY, -- UUID
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  last_fetched_at TIMESTAMP,
  name TEXT NOT NULL,
  url TEXT NOT NULL,
  user_id uuid NOT NULL REFERENCES users
                        ON DELETE CASCADE,
  UNIQUE(url),
  FOREIGN KEY(user_id) REFERENCES users (id)
);

-- +goose Down
DROP TABLE feeds;
