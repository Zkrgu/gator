-- +goose Up
CREATE TABLE feeds(
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	created_at TIMESTAMP NOT NULL DEFAULT now(),
	updated_at TIMESTAMP NOT NULL DEFAULT now(),
	name TEXT NOT NULL UNIQUE,
	url TEXT NOT NULL UNIQUE,
	user_id UUID NOT NULL references users(id) ON DELETE CASCADE
);

CREATE TRIGGER update_modified_feeds BEFORE UPDATE ON feeds FOR EACH ROW EXECUTE PROCEDURE update_modified_column();

-- +goose Down
DROP TABLE feeds;
