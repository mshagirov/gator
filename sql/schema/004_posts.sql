-- +goose Up
CREATE TABLE posts (
  id uuid PRIMARY KEY, -- UUID
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  title TEXT NOT NULL,
  url TEXT NOT NULL,
  description TEXT,
  published_at TIMESTAMP NOT NULL,
  feed_id uuid NOT NULL REFERENCES feeds
                        ON DELETE CASCADE,
  UNIQUE(url),
  FOREIGN KEY(feed_id) REFERENCES feeds (id)
);

-- +goose Down
DROP TABLE posts;
