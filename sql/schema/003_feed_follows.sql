-- +goose Up
CREATE TABLE feed_follows(
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	created_at TIMESTAMP NOT NULL DEFAULT now(),
	updated_at TIMESTAMP NOT NULL DEFAULT now(),
	user_id UUID NOT NULL references users(id) ON DELETE CASCADE,
	feed_id UUID NOT NULL references feeds(id) ON DELETE CASCADE,
	UNIQUE(user_id, feed_id)
);

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION create_feed_follows()
	RETURNS TRIGGER AS $$
	DECLARE
	BEGIN
		INSERT INTO feed_follows(id, user_id, feed_id) VALUES (gen_random_uuid(), NEW.user_id, NEW.id);
		RETURN NEW;
	END;
	$$ LANGUAGE 'plpgsql';
-- +goose StatementEnd

CREATE TRIGGER create_feed_follows_t AFTER INSERT ON feeds FOR EACH ROW EXECUTE PROCEDURE create_feed_follows();
CREATE TRIGGER update_modified_feeds BEFORE UPDATE ON feed_follows FOR EACH ROW EXECUTE PROCEDURE update_modified_column();

-- +goose Down
DROP TABLE feed_follows;

DROP TRIGGER create_feed_follows_t ON feeds;
