DROP DATABASE IF EXISTS go_takemikazuchi_api;
CREATE DATABASE go_takemikazuchi_api;
USE go_takemikazuchi_api;

SELECT * FROM workers;
SELECT * FROM jobs;
SELECT * FROM job_resources;
SELECT * FROM job_applications;
SELECT * FROM worker_resources;
SELECT * FROM worker_wallets;
SELECT * FROM user_addresses;
SELECT * FROM users;
SELECT * FROM withdrawals;
SELECT * FROM orders;
SELECT * FROM transactions;
SELECT * FROM categories;
UPDATE users SET role = 'Admin' WHERE id = 8;
SELECT * FROM workers;

DELETE FROM workers;
DELETE FROM worker_resources;
DELETE FROM worker_wallets;
DELETE FROM transactions;


DROP TABLE IF EXISTS reviews;

