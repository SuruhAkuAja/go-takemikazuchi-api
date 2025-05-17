CREATE TABLE workers
(
    id                     BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT NOT NULL,
    user_id                BIGINT UNSIGNED UNIQUE                     NOT NULL,
    rating                 FLOAT        DEFAULT 0,
    revenue                INT UNSIGNED DEFAULT 0,
    completed_jobs         INT UNSIGNED DEFAULT 0,
    location               VARCHAR(255),
    availability           BOOLEAN      DEFAULT TRUE,
    verified               BOOLEAN      DEFAULT FALSE,
    emergency_phone_number VARCHAR(30),
    created_at             TIMESTAMP    DEFAULT CURRENT_TIMESTAMP,
    updated_at             TIMESTAMP    DEFAULT CURRENT_TIMESTAMP,
    verified_at            TIMESTAMP    DEFAULT NULL,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);
