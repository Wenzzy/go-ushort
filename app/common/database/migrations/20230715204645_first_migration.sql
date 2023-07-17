-- +goose Up
-- +goose StatementBegin
CREATE TABLE "user"
(
    "id"            serial PRIMARY KEY NOT NULL,
    "email"         varchar(150)       NOT NULL UNIQUE,
    "password"      varchar(255)       NOT NULL,
    "role"          varchar(20)        NOT NULL DEFAULT 'user',
    "registered_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);


CREATE TABLE "link"
(
    "id"              serial PRIMARY KEY NOT NULL,
    "name"            varchar(100),
    "real_url"        TEXT               NOT NULL UNIQUE,
    "generated_alias" varchar(150)       NOT NULL UNIQUE,
    "created_at"      TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    "user_id"         integer,
    FOREIGN KEY (user_id) REFERENCES "user" (id)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE "link";
DROP TABLE "user";
-- +goose StatementEnd
