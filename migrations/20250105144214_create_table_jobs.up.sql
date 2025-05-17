CREATE TABLE jobs
(
    id          BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    user_id     BIGINT UNSIGNED         NOT NULL,
    address_id  BIGINT UNSIGNED         NOT NULL,
    category_id BIGINT UNSIGNED         NOT NULL,
    worker_id   BIGINT UNSIGNED         NOT NULL,
    title       VARCHAR(255)            NOT NULL,
    description TEXT                    NOT NULL,
    price       DECIMAL(10, 2) UNSIGNED NOT NULL,
    status      ENUM ('Open', 'Process', 'On Working','Closed', 'Done') DEFAULT 'Open',
    created_at  TIMESTAMP                                               DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP                                               DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (category_id) REFERENCES categories (id),
    FOREIGN KEY (user_id) REFERENCES users (id),
    FOREIGN KEY (address_id) REFERENCES user_addresses (id),
    FOREIGN KEY (worker_id) REFERENCES workers (id)
);
