-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
    `id` INT UNSIGNED auto_increment NOT NULL,
    `name` varchar(255) NOT NULL,
    `email` varchar(100) NOT NULL,
    `password` varchar(255) NOT NULL,
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    CONSTRAINT users_ID PRIMARY KEY (`id`),
    CONSTRAINT users_EMAIL UNIQUE KEY (`email`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8 COLLATE = utf8_general_ci;

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE users;

-- +goose StatementEnd