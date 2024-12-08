ALTER TABLE users ADD COLUMN email_confirmation_code VARCHAR(255) NULL AFTER password_hash;
ALTER TABLE users ADD COLUMN is_email_confirmed BOOLEAN DEFAULT FALSE AFTER email_confirmation_code;