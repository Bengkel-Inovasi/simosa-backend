-- Nodes table to store sensor node information
CREATE TABLE IF NOT EXISTS nodes (
    mac_address VARCHAR(17) PRIMARY KEY,
    alias VARCHAR(255),
    latitude DOUBLE PRECISION,
    longitude DOUBLE PRECISION,
    registered_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    is_registered BOOLEAN DEFAULT FALSE
);

-- Sensor readings table to store historical data from nodes
CREATE TABLE IF NOT EXISTS sensor_readings (
    id SERIAL PRIMARY KEY,
    node_mac VARCHAR(17) REFERENCES nodes(mac_address),
    ph REAL,
    n REAL,
    p REAL,
    k REAL,
    moisture REAL,
    temperature REAL,
    timestamp TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Harvests table to store harvest records
CREATE TABLE IF NOT EXISTS harvests (
    id SERIAL PRIMARY KEY,
    harvest_date DATE NOT NULL,
    yield_kg REAL NOT NULL,
    price_per_kg REAL NOT NULL,
    gross_income REAL NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Expenses table to store custom expenses related to a harvest
CREATE TABLE IF NOT EXISTS expenses (
    id SERIAL PRIMARY KEY,
    harvest_id INTEGER REFERENCES harvests(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    amount REAL NOT NULL
);
