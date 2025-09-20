-- Create device_gps table
CREATE TABLE device_gps (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    brand VARCHAR(100) NOT NULL,
    model VARCHAR(100) NOT NULL,
    imei VARCHAR(15) UNIQUE NOT NULL,
    protocol VARCHAR(20) NOT NULL,
    secret_key VARCHAR(255) NOT NULL,
    created_by INTEGER NOT NULL,
    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITHOUT TIME ZONE,
    deleted_at TIMESTAMP WITHOUT TIME ZONE,
    CONSTRAINT fk_device_gps_created_by FOREIGN KEY (created_by) REFERENCES customers(id)
);

-- Create indexes for device_gps
CREATE INDEX idx_device_gps_imei ON device_gps(imei);
CREATE INDEX idx_device_gps_created_by ON device_gps(created_by);
CREATE INDEX idx_device_gps_deleted_at ON device_gps(deleted_at);