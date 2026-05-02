CREATE TABLE channel_configs (
    id               UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    route_id         UUID         NOT NULL,
    channel_type     VARCHAR(20)  NOT NULL,
    to_field         VARCHAR(255),
    to_fixed         VARCHAR(255),
    subject          VARCHAR(500),
    template_id      UUID,
    message_template TEXT,
    webhook_url      VARCHAR(500),
    webhook_secret   VARCHAR(255),
    slack_channel    VARCHAR(255),
    active           BOOLEAN      NOT NULL DEFAULT true,
    created_at       TIMESTAMP    NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_channel_configs_route_id
        FOREIGN KEY (route_id) REFERENCES routes(id) ON DELETE CASCADE
);
