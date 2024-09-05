CREATE TABLE IF NOT EXISTS users (
    id UUID DEFAULT gen_random_uuid(),
    name text NOT NULL,
    email text NOT NULL UNIQUE,
	password text NOT NULL,
    create_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    update_time TIMESTAMP,
	CONSTRAINT pk PRIMARY KEY (id)
);