-- +goose Up
alter TABLE users 
add column hashed_password text not null default '????';

-- +goose Down
alter TABLE users 
drop column hashed_password;