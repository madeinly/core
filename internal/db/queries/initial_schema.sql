-- First statement: Create settings table
CREATE TABLE IF NOT EXISTS settings (
    name VARCHAR(255) PRIMARY KEY,
    value TEXT NOT NULL,
    type VARCHAR(50) NOT NULL DEFAULT 'string',
    package VARCHAR(255) NOT NULL,
    set_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    set_by VARCHAR(255),
    CONSTRAINT unique_package_option UNIQUE (package, name)
);

-- Second statement: Create settings_history table
CREATE TABLE IF NOT EXISTS settings_history (
    id SERIAL PRIMARY KEY,
    option_name VARCHAR(255) NOT NULL,
    old_value TEXT,
    new_value TEXT,
    changed_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    changed_by VARCHAR(255),
    FOREIGN KEY (option_name) REFERENCES settings(name) ON DELETE CASCADE
);

-- Third statement: Create migrations table
CREATE TABLE IF NOT EXISTS migrations (
    package TEXT PRIMARY KEY,
    executed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Fourth statement: Insert core migration
INSERT OR IGNORE INTO migrations (package) VALUES ('core');