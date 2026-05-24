CREATE TABLE routes (
    id                    UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id               UUID         NOT NULL,
    kafka_connection_id   UUID         NOT NULL,
    name                  VARCHAR(255) NOT NULL,
    topic                 VARCHAR(255) NOT NULL,
    active                BOOLEAN      NOT NULL DEFAULT true,
    created_at            TIMESTAMP    NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_routes_user_id
        FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,

    CONSTRAINT fk_routes_kafka_connection_id
        FOREIGN KEY (kafka_connection_id) REFERENCES kafka_connections(id) ON DELETE CASCADE
);
