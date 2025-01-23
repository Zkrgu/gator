-- name: CreateFeedFollowsForUser :one
INSERT INTO feed_follows (user_id, feed_id)
VALUES (
    $1,
    $2
)
RETURNING *;

-- name: GetFeedFollowsForUser :many
SELECT * FROM feed_follows JOIN users ON feed_follows.user_id = users.id JOIN feeds ON feed_follows.feed_id = feeds.id WHERE feed_follows.user_id = $1;

-- name: DeleteFeedFollowForUser :exec
DELETE FROM feed_follows USING feeds WHERE feed_follows.feed_id = feeds.id AND feed_follows.user_id = $1 AND feeds.url = $2;
