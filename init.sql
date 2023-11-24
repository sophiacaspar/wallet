

-- Create the database if it doesn't exist
DO $$ 
BEGIN
    IF NOT EXISTS (SELECT FROM pg_database WHERE datname = 'wallets_db') THEN
        CREATE DATABASE wallets_db;
    END IF;
END $$;

-- Switch to the created database
\c wallets_db;

-- Create the necessary table for wallets
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM   information_schema.tables 
        WHERE  table_schema = 'public'
        AND    table_name = 'wallets'
    ) THEN
        CREATE TABLE wallets (
            id bigserial PRIMARY KEY,
            playerId varchar(255) NOT NULL,
            balance integer NOT NULL DEFAULT 0,
            lastUpdatedAt timestamp(0) with time zone NOT NULL DEFAULT NOW(),
            UNIQUE (playerId)
        );
    END IF;
END $$;

-- Create a role (user) with a password
DO $$
BEGIN
    IF NOT EXISTS (SELECT FROM pg_catalog.pg_roles WHERE rolname = 'wallets') THEN
        CREATE ROLE wallets WITH LOGIN PASSWORD 'pa55w0rd';
    END IF;
END $$;

-- Grant specific permissions to user wallets
GRANT ALL PRIVILEGES ON DATABASE wallets_db TO wallets;
GRANT ALL PRIVILEGES ON TABLE wallets TO wallets;

-- Insert initial wallet data
INSERT INTO wallets (playerId, balance) VALUES
    ('player123', 100),
    ('player456', 200),
    ('player789', 300);