CREATE TABLE notification_logs (
    id                UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    route_id          UUID         NOT NULL,
    channel_config_id UUID         NOT NULL,
    topic             VARCHAR(255) NOT NULL,
    channel           VARCHAR(20)  NOT NULL,
    recipient         VARCHAR(255) NOT NULL,
    status            VARCHAR(20)  NOT NULL,
    payload           TEXT         NOT NULL,
    error_message     TEXT,
    attempts          INTEGER      NOT NULL DEFAULT 1,
    sent_at           TIMESTAMP,
    created_at        TIMESTAMP    NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_notification_logs_route_id
        FOREIGN KEY (route_id) REFERENCES routes(id) ON DELETE CASCADE,

    CONSTRAINT fk_notification_logs_channel_config_id
        FOREIGN KEY (channel_config_id) REFERENCES channel_configs(id) ON DELETE CASCADE
);
