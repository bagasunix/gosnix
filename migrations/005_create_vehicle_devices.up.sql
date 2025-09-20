-- Tabel vehicle_devices
CREATE TABLE vehicle_devices (
    id SERIAL PRIMARY KEY,
    vehicle_id UUID NOT NULL,
    device_id UUID NOT NULL,
    start_time TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    end_time TIMESTAMP WITHOUT TIME ZONE,
    is_active BOOLEAN DEFAULT TRUE NOT NULL,
    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITHOUT TIME ZONE,
    deleted_at TIMESTAMP WITHOUT TIME ZONE,
    CONSTRAINT fk_vehicle_devices_vehicle FOREIGN KEY (vehicle_id) REFERENCES vehicles(id),
    CONSTRAINT fk_vehicle_devices_device FOREIGN KEY (device_id) REFERENCES device_gps(id)
);

-- Index untuk vehicle_devices
CREATE INDEX idx_vehicle_devices_vehicle_id ON vehicle_devices(vehicle_id);
CREATE INDEX idx_vehicle_devices_device_id ON vehicle_devices(device_id);
CREATE INDEX idx_vehicle_devices_is_active ON vehicle_devices(is_active);
CREATE INDEX idx_vehicle_devices_deleted_at ON vehicle_devices(deleted_at);
CREATE INDEX idx_vehicle_devices_time_range ON vehicle_devices(start_time, end_time);