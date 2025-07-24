-- +goose Up
-- +goose StatementBegin
-- +goose ENVSUB ON
CREATE USER ${AUTH_USER} WITH PASSWORD '${AUTH_PASS}';
CREATE SCHEMA auth AUTHORIZATION ${AUTH_USER};
ALTER ROLE ${AUTH_USER} SET search_path TO auth;
GRANT USAGE ON SCHEMA auth TO ${AUTH_USER};
GRANT CREATE ON SCHEMA auth TO ${AUTH_USER};
-- +goose ENVSUB OFF
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- +goose ENVSUB ON
DROP USER IF EXISTS ${AUTH_USER};
DROP SCHEMA IF EXISTS auth CASCADE;
-- +goose ENVSUB OFF
-- +goose StatementEnd
