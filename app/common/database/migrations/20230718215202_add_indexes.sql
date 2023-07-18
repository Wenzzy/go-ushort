-- +goose Up
-- +goose StatementBegin
CREATE INDEX "link_real_url_idx" ON "link" ("real_url");
CREATE INDEX "link_generated_alias_idx" ON "link" ("generated_alias");
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX "link_real_url_idx";
DROP INDEX "link_generated_alias_idx";
-- +goose StatementEnd
