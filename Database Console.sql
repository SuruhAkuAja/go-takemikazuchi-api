DROP DATABASE IF EXISTS go_takemikazuchi_api;
CREATE DATABASE go_takemikazuchi_api;
USE go_takemikazuchi_api;

SELECT *
FROM workers;
SELECT *
FROM jobs;
SELECT *
FROM job_resources;
SELECT *
FROM job_applications;
SELECT *
FROM worker_resources;
SELECT *
FROM workers;
SELECT *
FROM worker_wallets;
SELECT *
FROM user_addresses;
SELECT *
FROM users;
SELECT *
FROM withdrawals;
SELECT *
FROM orders;
SELECT *
FROM transactions;
SELECT *
FROM categories;
UPDATE users
SET role = 'Admin'
WHERE id = 8;
SELECT *
FROM workers;
SELECT *
FROM users;

DELETE
FROM workers;
DELETE
FROM worker_resources;
DELETE
FROM worker_wallets;
DELETE
FROM transactions;


DROP TABLE IF EXISTS reviews;

SELECT *
FROM users
         LEFT JOIN users ON users.id = jobs.user_id
         LEFT JOIN categories ON categories.id = jobs.category_id
         LEFT JOIN workers ON workers.id = jobs.worker_id
         LEFT JOIN user_addresses ON user_addresses.id = jobs.address_id
WHERE jobs.id = 7;

SELECT *
FROM worker_wallets;