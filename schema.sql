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
	title STRING,
	content STRING,
	background_color STRING(7),
    create_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    update_time TIMESTAMP,
    delete_time TIMESTAMP,
	CONSTRAINT pk PRIMARY KEY (id),
	CONSTRAINT fk_user
    	FOREIGN KEY (user_id) 
    	REFERENCES users(id)
    	ON DELETE CASCADE
);