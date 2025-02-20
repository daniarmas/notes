CREATE TABLE IF NOT EXISTS users (
    id UUID DEFAULT gen_random_uuid(),
    name text NOT NULL,
    email text NOT NULL UNIQUE,
	password text NOT NULL,
    create_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    update_time TIMESTAMP,
	CONSTRAINT pk PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS notes (
	id UUID DEFAULT gen_random_uuid(),
	user_id UUID NOT NULL,
	title VARCHAR,
	content VARCHAR,
    create_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    update_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    delete_time TIMESTAMP,
	CONSTRAINT pk PRIMARY KEY (id),
	CONSTRAINT fk_user
    	FOREIGN KEY (user_id) 
    	REFERENCES users(id)
    	ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS files (
	id UUID DEFAULT gen_random_uuid(),
	processed_file VARCHAR,
	original_file VARCHAR NOT NULL,
	note_id UUID NOT NULL,
	create_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
	update_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
	delete_time TIMESTAMP,
	CONSTRAINT pk PRIMARY KEY (id),
	CONSTRAINT fk_note
		FOREIGN KEY (note_id) 
		REFERENCES notes(id)
		ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS refresh_tokens (
	id UUID DEFAULT gen_random_uuid(),
	user_id UUID NOT NULL,
	create_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
	update_time TIMESTAMP,
	CONSTRAINT pk PRIMARY KEY (id),
	CONSTRAINT fk_user
		FOREIGN KEY (user_id) 
		REFERENCES users(id)
		ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS access_tokens (
	id UUID DEFAULT gen_random_uuid(),
	user_id UUID NOT NULL,
	refresh_token_id UUID NOT NULL,
    create_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    update_time TIMESTAMP,
	CONSTRAINT pk PRIMARY KEY (id),
	CONSTRAINT fk_user
    	FOREIGN KEY (user_id) 
    	REFERENCES users(id)
    	ON DELETE CASCADE,
	CONSTRAINT fk_refresh_token
    	FOREIGN KEY (refresh_token_id) 
    	REFERENCES refresh_tokens(id)
    	ON DELETE CASCADE
);