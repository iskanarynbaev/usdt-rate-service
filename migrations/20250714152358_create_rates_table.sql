-- +goose Up
CREATE TABLE rates
(
    id        SERIAL PRIMARY KEY,
    ask       NUMERIC     NOT NULL,
    bid       NUMERIC     NOT NULL,
    timestamp TIMESTAMPTZ NOT NULL
);

-- +goose Down
DROP TABLE rates;
