CREATE TABLE templates (
    id         UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id    UUID         NOT NULL,
    name       VARCHAR(255) NOT NULL,
    content    TEXT         NOT NULL,
    created_at TIMESTAMP    NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP    NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_templates_user_id
        FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

ALTER TABLE channel_configs
    ADD CONSTRAINT fk_channel_configs_template_id
        FOREIGN KEY (template_id) REFERENCES templates(id) ON DELETE SET NULL;
