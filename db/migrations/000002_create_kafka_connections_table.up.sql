CREATE TABLE kafka_connections (
    id         UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id    UUID         NOT NULL,
    name       VARCHAR(255) NOT NULL,
    brokers    VARCHAR(500) NOT NULL,
    group_id   VARCHAR(255) NOT NULL,
    active     BOOLEAN      NOT NULL DEFAULT true,
    created_at TIMESTAMP    NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_kafka_connections_user_id
        FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
