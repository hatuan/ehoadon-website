-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied

-- +goose StatementBegin
DO $$DECLARE
    _CurrencyLCYId bigint;
    _UserID bigint;
    _OrganizationID bigint;
    _ClientID bigint;
BEGIN

INSERT INTO "user_profile"("password_question", 
    "password_answer", 
    "password", 
    "salt", 
    "name", 
    "last_modified_date", 
    "last_login_ip", 
    "last_login_date", 
    "last_locked_out_reason", 
    "last_locked_out_date", 
    "is_locked_out", 
    "is_activated", 
    "full_name", 
    "email", 
    "created_date", 
    "comment", 
    "client_id", 
    "culture_ui_id") 
VALUES (
    '', 
    '', 
    '{{.UserProfile.Password}}', 
    '{{.UserProfile.Salt}}',  
    'ADMIN', 
    NOW(), 
    '0.0.0.0', 
    NOW(), 
    '', 
    NOW(), 
    'f', 
    't', 
    'Administrator', 
    '{{.UserProfile.Email}}', 
    '2014-04-25 10:22:39.007887+07', 
    '', 
    {{.Client.ClientID}},
    'vi-VN') RETURNING id INTO _UserID;

INSERT INTO "client"(
    "version", 
    "description", 
    "is_activated", 
    "id", 
    "culture_id", 
    "amount_decimal_places", 
    "amount_rounding_precision", 
    "unit_amount_decimal_places", 
    "unit_amount_rounding_precision", 
    "currency_lcy_id",
    "vat_number",
    "group_unit_code",
    "vat_method_code",
    "province_code",
    "districts_code",
    "address",
    "address_transition",
    "telephone",
    "email",
    "fax",
    "website",
    "representative_name",
    "representative_position",
    "contact_name",
    "mobile",
    "bank_account",
    "bank_name",
    "tax_authorities_id", 
    "rec_created_by", 
    "rec_modified_by", 
    "rec_created_at", 
    "rec_modified_at") 
VALUES ( 
    1, 
    '{{.Client.Description}}', 
    't', 
    {{.Client.ClientID}}, 
    'vi-VN', 
    3, 
    0.001, 
    3, 
    0.001, 
    0, --currency_lcy_id will update at last
    '{{.Client.VatNumber}}',
    '',
    '',
    '',
    '',
    '{{.Client.Address}}',
    '{{.Client.AddressTransition}}',
    '{{.Client.Telephone}}',
    '{{.Client.Email}}',
    '{{.Client.Fax}}',
    '{{.Client.Website}}',
    '{{.Client.RepresentativeName}}',
    '{{.Client.RepresentativePosition}}',
    '{{.Client.ContactName}}',
    '{{.Client.Mobile}}',
    '',
    '',
    0,
    _UserID, 
    _UserID, 
    NOW(), 
    NOW()) RETURNING id INTO _ClientID;

INSERT INTO "organization"(
    "version", 
    "status", 
    "rec_modified_by", 
    "rec_created_by", 
    "description", 
    "code", 
    "client_id", 
    "rec_created_at", 
    "rec_modified_at") 
VALUES (1, 
    1, 
    _UserID,
    _UserID, 
    'All Organization', 
    '*', 
    _ClientID, 
    NOW(), 
    NOW()
) RETURNING id INTO _OrganizationID ;


INSERT INTO "currency"(
    "version", 
    "status", 
    "rec_modified_by", 
    "rec_created_by", 
    "organization_id", 
    "description", 
    "code", 
    "client_id", 
    "rec_created_at", 
    "rec_modified_at") 
VALUES (
    1, 
    1, 
    _UserID, 
    _UserID, 
    _OrganizationID, 
    'Việt Nam Đồng', 
    'VND', 
    _ClientID, 
    NOW(), 
    NOW()) RETURNING id INTO _CurrencyLCYId;

UPDATE client SET currency_lcy_id = _CurrencyLCYId;

END$$;
-- +goose StatementEnd

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DELETE FROM "organization";
DELETE FROM "client";
DELETE FROM "user_profile";
DELETE FROM "currency";