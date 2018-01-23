
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE SEQUENCE IF NOT EXISTS global_id_sequence;

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION id_generator(OUT result bigint) AS
$BODY$
DECLARE
    --https://engineering.instagram.com/sharding-ids-at-instagram-1cf5a71e5a5c#.1entthc2d
    --http://rob.conery.io/2014/05/29/a-better-id-generator-for-postgresql/
    --https://www.depesz.com/2015/10/30/is-c-faster-for-instagram-style-id-generation/
    --shard_id : 5 ... 2000
    our_epoch bigint := 1314220021721;
    seq_id bigint;
    now_millis bigint;
    shard_id int := 5;
BEGIN
    SELECT nextval('global_id_sequence') % 1024 INTO seq_id;

    SELECT FLOOR(EXTRACT(EPOCH FROM clock_timestamp()) * 1000) INTO now_millis;
    result := (now_millis - our_epoch) << 23; -- <<(64-41)
    result := result | (shard_id <<10); --<<(64-41-13)
    result := result | (seq_id);
END

$BODY$
  LANGUAGE plpgsql;    
-- +goose StatementEnd


CREATE TABLE IF NOT EXISTS client
(
  id bigint NOT NULL DEFAULT id_generator(),
  description character varying NOT NULL,
  is_activated boolean NOT NULL,
  activated_code character varying NOT NULL,
  culture_id character varying NOT NULL,
  amount_decimal_places smallint NOT NULL,
  amount_rounding_precision numeric NOT NULL,
  unit_amount_decimal_places smallint NOT NULL,
  unit_amount_rounding_precision numeric NOT NULL,
  currency_lcy_id bigint,
  vat_number character varying NOT NULL, --ma so thue
  group_unit_code character varying NOT NULL, -- khoi doanh nghiep : doanh nghiep, truong hoc, y te, khac
  vat_method_code character varying NOT NULL, -- Hình thức khai thuế : Trực Tiếp, Khấu trừ, Đơn vị trong khu chế xuất, khu phi thuế quan
  province_code character varying NOT NULL, -- Tỉnh / Thành phố
  districts_code character varying NOT NULL, -- Quan / Huyen
  address character varying NOT NULL, -- Dia chi
  address_transition character varying NOT NULL, -- Dia chi giao dich
  telephone character varying NOT NULL, -- Dien thoai
  email character varying NOT NULL, -- Email
  fax character varying NOT NULL, -- FAX
  website character varying NOT NULL, -- website
  representative_name character varying NOT NULL, -- nguoi dai dien
  representative_position character varying NOT NULL, -- chuc vu nguoi dai dien
  contact_name character varying NOT NULL, -- nguoi lien he
  mobile character varying NOT NULL, -- mobile
  bank_account character varying NOT NULL, -- tai khoan ngan hang
  bank_name character varying NOT NULL, -- ten ngan hang
  tax_authorities_id bigint, -- Chi cuc thue tax_authorities
  version bigint NOT NULL,
  rec_created_by bigint NOT NULL,
  rec_modified_by bigint NOT NULL,
  rec_created_at timestamp with time zone NOT NULL,
  rec_modified_at timestamp with time zone NOT NULL,
  CONSTRAINT pk_client PRIMARY KEY (id)
);

CREATE INDEX IF NOT EXISTS idx_client_currency_lcy ON client USING btree (currency_lcy_id);

CREATE UNIQUE INDEX IF NOT EXISTS idx_client_active_code ON client USING btree (activated_code);

CREATE INDEX IF NOT EXISTS idx_client_vat_number ON client USING btree (vat_number);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE IF EXISTS client;
DROP FUNCTION IF EXISTS id_generator;
DROP SEQUENCE IF EXISTS global_id_sequence;
