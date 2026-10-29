-- +goose Up
CREATE TABLE rate_updates (
    id              UUID            PRIMARY KEY,
    base_currency   TEXT            NOT NULL,
    quote_currency  TEXT            NOT NULL,
    status          TEXT            NOT NULL,
    price           NUMERIC(20, 10),
    rate_timestamp  TIMESTAMPTZ,
    error           TEXT,
    attempts        INT             NOT NULL DEFAULT 0,
    next_attempt_at TIMESTAMPTZ     NOT NULL DEFAULT now(),
    idempotency_key TEXT            UNIQUE,
    created_at      TIMESTAMPTZ     NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ     NOT NULL DEFAULT now()
);

CREATE INDEX rate_updates_pending_idx ON rate_updates (next_attempt_at) WHERE status = 'pending';

CREATE TABLE rates (
    base_currency  TEXT            NOT NULL,
    quote_currency TEXT            NOT NULL,
    price          NUMERIC(20, 10) NOT NULL,
    rate_timestamp TIMESTAMPTZ     NOT NULL,
    updated_at     TIMESTAMPTZ     NOT NULL,
    PRIMARY KEY (base_currency, quote_currency)
);

-- +goose Down
DROP TABLE rates;
DROP TABLE rate_updates;
