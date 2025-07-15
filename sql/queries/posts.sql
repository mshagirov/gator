-- name: CreatePost :one
INSERT INTO posts (
    id,
    created_at,
    updated_at,
    title,
    url,
    description,
    published_at,
    feed_id
)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7,
    $8
)
RETURNING *;
--

-- name: GetPostsForUserId :many
WITH user_feed_ids AS (
  SELECT
    feed_id
  FROM feed_follows
  WHERE user_id = $1
)
SELECT
  posts.*
FROM user_feed_ids
INNER JOIN posts ON user_feed_ids.feed_id = posts.feed_id
ORDER BY posts.published_at DESC NULLS LAST
LIMIT $2 
;
--
