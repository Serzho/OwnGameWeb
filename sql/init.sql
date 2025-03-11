CREATE TABLE "user"(
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    password CHAR(60) NOT NULL
);

-- CONSTRAINTS FOR USER TABLE
ALTER TABLE "user"
    ADD CONSTRAINT user_name_length CHECK (LENGTH(name) >= 3),
    ADD CONSTRAINT user_email_unique UNIQUE (email),
    ADD CONSTRAINT user_email_format CHECK (
        email ~* '^[a-z0-9._+%-]+@[a-z0-9.-]+\.[a-z]{2,4}$'
        ),
    ADD CONSTRAINT user_password_length CHECK (LENGTH(password) = 60);

CREATE UNIQUE INDEX user_email_unique_index
    ON "user" (email);


CREATE TABLE "game"(
    id SERIAL PRIMARY KEY,
    title VARCHAR(64) NOT NULL,
    status VARCHAR(16) NOT NULL,
    invite_code CHAR(6),
    start_time timestamp NOT NULL,
    master_id INT REFERENCES "user"(id),
    players_ids INT[] NOT NULL,
    max_players SMALLINT NOT NULL
);

-- CONSTRAINTS FOR GAME TABLE
ALTER TABLE "game"
    ADD CONSTRAINT game_title_length CHECK (LENGTH(title) >= 5),
    ADD CONSTRAINT game_status_valid CHECK (status ~ 'created|inprocess|firststage|secondstage|thirdstage|finished|archieved'),
    ADD CONSTRAINT game_invite_code_valid CHECK (invite_code ~ '[a-zA-Z0-9]+'),
    ADD CONSTRAINT game_max_users_valid CHECK (max_players > 1 AND max_players <= 6),
    ADD CONSTRAINT game_users_count CHECK (array_length(players_ids, 1) <= max_players);






