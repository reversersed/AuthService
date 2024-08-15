CREATE TABLE IF NOT EXISTS sessions (
    id uuid PRIMARY KEY NOT NULL UNIQUE DEFAULT gen_random_uuid(),
    clientid varchar(32) not null,
    refreshtoken bytea not null unique
);