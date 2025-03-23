CREATE OR REPLACE FUNCTION is_array_unique(arr integer[])
    RETURNS boolean AS $$
BEGIN
    IF arr IS NULL OR array_length(arr, 1) IS NULL THEN
        RETURN true;
    END IF;

    RETURN (
        SELECT COUNT(*) = COUNT(DISTINCT x)
        FROM unnest(arr) AS t(x)
    );
END;
$$ LANGUAGE plpgsql IMMUTABLE;



CREATE TABLE "user"(
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    password CHAR(60) NOT NULL,
    packs int[] NOT NULL DEFAULT '{}',
    played_games INT NOT NULL DEFAULT 0,
    won_games INT NOT NULL DEFAULT 0
);

-- CONSTRAINTS FOR USER TABLE
ALTER TABLE "user"
    ADD CONSTRAINT user_name_length CHECK (LENGTH(name) >= 3),
    ADD CONSTRAINT user_email_unique UNIQUE (email),
    ADD CONSTRAINT user_email_format CHECK (
        email ~* '^[a-z0-9._+%-]+@[a-z0-9.-]+\.[a-z]{2,4}$'
        ),
    ADD CONSTRAINT user_packs_valid CHECK(is_array_unique(packs)),
    ADD CONSTRAINT user_password_length CHECK (LENGTH(password) = 60),
    ADD CONSTRAINT user_won_games_valid CHECK(won_games <= played_games);

CREATE UNIQUE INDEX user_email_unique_index
    ON "user" (email);



CREATE TABLE "question_pack"(
                                id SERIAL PRIMARY KEY,
                                title VARCHAR(32) NOT NULL,
                                filename VARCHAR(64) NOT NULL,
                                owner INT REFERENCES "user"(id)
);

ALTER TABLE "question_pack"
    ADD CONSTRAINT question_pack_filename_length CHECK (LENGTH(filename) = 64),
    ADD CONSTRAINT question_pack_title_length CHECK (LENGTH(title) > 6);




CREATE TABLE "question_sample"(
                                  id SERIAL PRIMARY KEY,
                                  pack INT REFERENCES question_pack(id),
                                  content json NOT NULL
);

CREATE TABLE "game"(
                       id SERIAL PRIMARY KEY,
                       title VARCHAR(64) NOT NULL,
                       status VARCHAR(16) NOT NULL,
                       invite_code CHAR(6),
                       start_time timestamp NOT NULL,
                       master INT REFERENCES "user"(id),
                       players_ids INT[] NOT NULL DEFAULT '{}',
                       max_players SMALLINT NOT NULL,
                       sample INT REFERENCES "question_sample"(id),
                       used_questions INT[] NOT NULL DEFAULT '{}',
                       game_score JSON NOT NULL,
                       winner INT REFERENCES "user"(id)
);

-- CONSTRAINTS FOR GAME TABLE
ALTER TABLE "game"
    ADD CONSTRAINT game_title_length CHECK (LENGTH(title) >= 5),
    ADD CONSTRAINT game_status_valid CHECK (status ~ 'created|firststage|secondstage|thirdstage|finished|archieved'),
    ADD CONSTRAINT game_invite_code_valid CHECK (invite_code ~ '[a-zA-Z0-9]+'),
    ADD CONSTRAINT game_max_users_valid CHECK (max_players > 1 AND max_players <= 6),
    ADD CONSTRAINT game_players_valid CHECK(is_array_unique(players_ids)),
    ADD CONSTRAINT game_used_questions_valid CHECK(is_array_unique(used_questions)),
    ADD CONSTRAINT game_users_count CHECK (array_length(players_ids, 1) <= max_players);





