-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS materials (
  uuid UUID PRIMARY KEY,
  type_materials VARCHAR(50) NOT NULL,
  status VARCHAR(50) NOT NULL,
  title VARCHAR(255) NOT NULL,
  content TEXT NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS materials;
-- +goose StatementEnd
