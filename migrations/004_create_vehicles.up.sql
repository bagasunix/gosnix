-- Tabel vehicles
CREATE TABLE vehicles (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    category_id UUID NOT NULL,
    customer_id INTEGER NOT NULL,
    plate_no VARCHAR(20) UNIQUE NOT NULL,
    model VARCHAR(100) NOT NULL,
    brand VARCHAR(100) NOT NULL,
    manufacture_year INTEGER NOT NULL,
    color VARCHAR(50) NOT NULL,
    vin VARCHAR(50) UNIQUE NOT NULL,
    fuel_type VARCHAR(50) NOT NULL,
    engine_capacity DECIMAL(6,2) NOT NULL,
    max_speed INTEGER NOT NULL,
    is_active INTEGER NOT NULL DEFAULT 1,
    created_by INTEGER NOT NULL,
    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITHOUT TIME ZONE,
    deleted_at TIMESTAMP WITHOUT TIME ZONE,
    CONSTRAINT fk_vehicles_category FOREIGN KEY (category_id) REFERENCES vehicle_categories(id),
    CONSTRAINT fk_vehicles_customer FOREIGN KEY (customer_id) REFERENCES customers(id)
);

-- Index untuk vehicles
CREATE INDEX idx_vehicles_plate_no ON vehicles(plate_no);
CREATE INDEX idx_vehicles_vin ON vehicles(vin);
CREATE INDEX idx_vehicles_category_id ON vehicles(category_id);
CREATE INDEX idx_vehicles_customer_id ON vehicles(customer_id);
CREATE INDEX idx_vehicles_created_by ON vehicles(created_by);
CREATE INDEX idx_vehicles_deleted_at ON vehicles(deleted_at);
CREATE INDEX idx_vehicles_is_active ON vehicles(is_active);

ALTER TABLE vehicles ADD CONSTRAINT fk_vehicles_created_by FOREIGN KEY (created_by) REFERENCES customers(id);