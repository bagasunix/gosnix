-- Create vehicle_categories table
CREATE TABLE vehicle_categories (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) UNIQUE NOT NULL,
    description TEXT,
    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITHOUT TIME ZONE,
    deleted_at TIMESTAMP WITHOUT TIME ZONE
);

-- Create indexes for vehicle_categories
CREATE INDEX idx_vehicle_categories_name ON vehicle_categories(name);
CREATE INDEX idx_vehicle_categories_deleted_at ON vehicle_categories(deleted_at);