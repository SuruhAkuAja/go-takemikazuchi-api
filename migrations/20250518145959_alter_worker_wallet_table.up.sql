ALTER TABLE worker_wallets
    ADD COLUMN balance DECIMAL(10, 2) NOT NULL DEFAULT 0.00;