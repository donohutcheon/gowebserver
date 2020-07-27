DROP DATABASE banking;

CREATE DATABASE banking;
\c banking

CREATE OR REPLACE FUNCTION trigger_set_timestamp()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = NOW();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TABLE users (
  id BIGSERIAL PRIMARY KEY,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  deleted_at TIMESTAMPTZ,
  logged_out_at TIMESTAMPTZ,
  email VARCHAR(255),
  password VARCHAR(255),
  role VARCHAR(255),
  state VARCHAR(16)
);

CREATE TRIGGER user_updated
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();

CREATE INDEX idx_users_email
ON users(email);

CREATE INDEX idx_users_deleted_at
ON users(deleted_at);

CREATE TABLE sign_up_confirmations (
  id BIGSERIAL PRIMARY KEY,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  deleted_at TIMESTAMPTZ,
  nonce VARCHAR(32) UNIQUE NOT NULL,
  user_id BIGINT NOT NULL,
  FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE
);

CREATE TRIGGER user_updated
BEFORE UPDATE ON sign_up_confirmations
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();

CREATE INDEX idx_sign_up_confirmations_user_id
ON sign_up_confirmations(user_id);

CREATE INDEX idx_sign_up_confirmations_nonce
ON sign_up_confirmations(nonce);

CREATE INDEX idx_sign_up_confirmations_deleted_at
ON sign_up_confirmations(deleted_at);

CREATE TABLE card_transactions (
  id BIGSERIAL PRIMARY KEY,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  datetime TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  amount BIGINT NOT NULL,
  currency_scale SMALLINT NOT NULL,
  currency_code VARCHAR(255) NOT NULL,
  reference VARCHAR(255) NOT NULL,
  merchant_name VARCHAR(255) NOT NULL,
  merchant_city VARCHAR(255) NOT NULL,
  merchant_country_code VARCHAR(255) NOT NULL,
  merchant_country_name VARCHAR(255) NOT NULL,
  merchant_category_code VARCHAR(255) NOT NULL,
  merchant_category_name VARCHAR(255) NOT NULL,
  user_id BIGINT NOT NULL,
  FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE
);

CREATE INDEX idx_card_transactions_user_id
ON card_transactions(user_id);

