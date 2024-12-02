CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    user_guid varchar(255),
    refresh_hash varchar(255),
    ip varchar(50)
)