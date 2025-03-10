CREATE TABLE "user"(
   id SERIAL PRIMARY KEY,
   name VARCHAR(255) NOT NULL,
   email VARCHAR(255) NOT NULL,
   password CHAR(60) NOT NULL
);

ALTER TABLE "user"
    ADD CONSTRAINT user_name_length CHECK (LENGTH(name) >= 3),
    ADD CONSTRAINT user_email_unique UNIQUE (email),
    ADD CONSTRAINT user_email_format CHECK (
        email ~* '^[a-z0-9._+%-]+@[a-z0-9.-]+\.[a-z]{2,4}$'
        ),
    ADD CONSTRAINT user_password_length CHECK (LENGTH(password) = 60);

CREATE UNIQUE INDEX user_email_unique_index
    ON "user" (email);