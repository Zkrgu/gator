-- +goose Up
CREATE TABLE posts(
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	created_at TIMESTAMP NOT NULL DEFAULT now(),
	updated_at TIMESTAMP NOT NULL DEFAULT now(),
	title TEXT NOT NULL,
	url TEXT UNIQUE NOT NULL,
	description TEXT NOT NULL,
	published_at TIMESTAMP NOT NULL,
	feed_id UUID NOT NULL references feeds(id) ON DELETE CASCADE
);

CREATE TRIGGER update_modified_feeds BEFORE UPDATE ON posts FOR EACH ROW EXECUTE PROCEDURE update_modified_column();

-- +goose Down
DROP TABLE posts;
