-- +goose Up
-- +goose StatementBegin
-- +goose ENVSUB ON
CREATE USER ${PROFILES_USER} WITH PASSWORD '${PROFILES_PASS}';

CREATE SCHEMA profiles AUTHORIZATION ${PROFILES_USER};

ALTER ROLE ${PROFILES_USER} SET search_path TO profiles;

GRANT USAGE ON SCHEMA profiles TO ${PROFILES_USER};

GRANT CREATE ON SCHEMA profiles TO ${PROFILES_USER};
-- +goose ENVSUB OFF
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- +goose ENVSUB ON
DROP USER IF EXISTS ${PROFILES_USER};

DROP SCHEMA IF EXISTS profiles CASCADE;
-- +goose ENVSUB OFF
-- +goose StatementEnd