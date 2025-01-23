-- name: CreateFeed :one
INSERT INTO feeds (name, url, user_id)
VALUES (
    $1,
    $2,
	$3
)
RETURNING *;

-- name: GetFeedByUrl :one
SELECT * FROM feeds WHERE url = $1;

-- name: GetFeeds :many
SELECT * FROM feeds JOIN users ON users.id = feeds.user_id;

-- name: MarkFeedFetched :exec
UPDATE feeds SET last_fetched_at = now() WHERE id = $1;

-- name: GetNextFeedToFetch :one
SELECT * FROM feeds ORDER BY last_fetched_at ASC NULLS FIRST LIMIT 1;
