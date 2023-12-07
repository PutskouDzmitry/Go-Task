BEGIN;

CREATE TABLE transactional (
    id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    wallet_id UUID NOT NULL REFERENCES wallet(id) on delete cascade,
    transactional_type VARCHAR(2) NOT NULL,
    amount VARCHAR(50) NOT NULL,
    updated_balance VARCHAR(50) NOT NULL
);

CREATE INDEX idx_transactional ON transactional(id);

COMMIT;