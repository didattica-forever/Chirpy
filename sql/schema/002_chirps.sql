-- +goose Up
CREATE TABLE chirps (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE not null  DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE not null  DEFAULT CURRENT_TIMESTAMP,
    body VARCHAR(255) NOT NULL,
    user_id UUID not null REFERENCES users(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE chirps;