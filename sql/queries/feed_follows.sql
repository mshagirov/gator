-- name: CreateFeedFollow :one
WITH inserted_feed_follow AS (
  INSERT INTO feed_follows (
    id,
    created_at,
    updated_at,
    user_id,
    feed_id
  )
  VALUES (
    $1,
    $2,
    $3,
    $4,
    $5
  )
  RETURNING *
)
SELECT
    inserted_feed_follow.*,
    feeds.name AS feed_name,
    users.name AS user_name
FROM inserted_feed_follow 
INNER JOIN feeds ON inserted_feed_follow.feed_id = feeds.id
INNER JOIN users ON inserted_feed_follow.user_id = users.id
;
--

-- name: GetFeedFollowsForUser :many
WITH user_feed_ids AS (
  SELECT
    feed_id,
    user_id
  FROM feed_follows 
  WHERE user_id=( SELECT id FROM users WHERE users.name= $1 )
)
SELECT
  users.name AS user_name,
  feeds.name AS feed_name,
  feeds.url AS feed_url
FROM user_feed_ids
INNER JOIN feeds ON user_feed_ids.feed_id = feeds.id
INNER JOIN users ON user_feed_ids.user_id = users.id
;
--

-- name: GetFeedFollowsForUserId :many
SELECT
  feed_follows.*,
  feeds.name AS feed_name,
  users.name AS user_name
FROM feed_follows
INNER JOIN feeds ON feed_follows.feed_id = feeds.id
INNER JOIN users ON feed_follows.user_id = users.id
WHERE feed_follows.user_id = $1
;
--
