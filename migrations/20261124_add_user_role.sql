ALTER TABLE users
    ADD COLUMN user_role VARCHAR(16) NOT NULL DEFAULT 'user';


UPDATE users SET user_role = 'admin' WHERE email = 'ingwilliamsespinosalopez@gmail.com';
