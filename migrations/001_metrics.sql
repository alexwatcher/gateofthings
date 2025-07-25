-- +goose Up
-- +goose StatementBegin
-- +goose ENVSUB ON
CREATE USER ${METRICS_USER} WITH PASSWORD '${METRICS_PASS}';
GRANT pg_monitor TO ${METRICS_USER};
-- +goose ENVSUB OFF
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- +goose ENVSUB ON
REVOKE pg_monitor FROM ${METRICS_USER};
DROP USER ${METRICS_USER};
-- +goose ENVSUB OFF
-- +goose StatementEnd
