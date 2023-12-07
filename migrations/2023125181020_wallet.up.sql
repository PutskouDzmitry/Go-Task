BEGIN;

CREATE TABLE wallet (
    id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    coin VARCHAR(20) NOT NULL,
    money VARCHAR(50) NOT NULL
);

CREATE INDEX idx_wallet ON wallet(id);

COMMIT;