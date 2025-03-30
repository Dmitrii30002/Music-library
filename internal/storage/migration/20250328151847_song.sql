-- +goose Up
-- +goose StatementBegin
CREATE TABLE song (
    id            SERIAL PRIMARY KEY,     -- идентификатор песни
    group         VARCHAR(255) NOT NULL,  -- название группа
    name          VARCHAR(255) NOT NULL,  -- название песни
    release_date  DATE, NOT NULL,          -- дата релиза
    text          TEXT,                   -- текст песни
    link          VARCHAR(2083)           -- ссылка на песню
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS song;
-- +goose StatementEnd
