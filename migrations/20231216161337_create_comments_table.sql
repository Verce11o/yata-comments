-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS comments(
    comment_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tweet_id UUID NOT NULL,
    user_id UUID NOT NULL ,
    text varchar(255) NOT NULL,
    image_name varchar(255) null,
    image_temp_url text null ,
    created_at   TIMESTAMP WITH TIME ZONE    NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMP WITH TIME ZONE             DEFAULT CURRENT_TIMESTAMP

);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
