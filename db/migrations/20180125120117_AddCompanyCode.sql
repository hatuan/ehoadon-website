
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE client ADD code character varying;
UPDATE client SET code = '';
ALTER TABLE client ALTER COLUMN code SET NOT NULL;

CREATE UNIQUE INDEX IF NOT EXISTS idx_client_code ON client USING btree (code);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
ALTER TABLE client DROP COLUMN code;
