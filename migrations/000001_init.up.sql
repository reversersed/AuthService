CREATE TABLE IF NOT EXISTS sessions (
    id uuid PRIMARY KEY NOT NULL UNIQUE DEFAULT gen_random_uuid(),
    clientid uuid not null,
    refreshtoken varchar(72) not null unique,
    created bigint not null,
    refreshed timestamp default null
);