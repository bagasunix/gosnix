-- Create location_updates table
CREATE TABLE location_updates (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    vehicle_id UUID NOT NULL,
    session_id UUID NOT NULL,
    latitude DECIMAL(10,8) NOT NULL,
    longitude DECIMAL(11,8) NOT NULL,
    speed DECIMAL(6,2),
    heading SMALLINT,
    received_at TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_location_updates_vehicle FOREIGN KEY (vehicle_id) REFERENCES vehicles(id),
    CONSTRAINT fk_location_updates_session FOREIGN KEY (session_id) REFERENCES tracking_sessions(id)
);

-- Create indexes for location_updates
CREATE INDEX idx_location_updates_vehicle_id ON location_updates(vehicle_id);
CREATE INDEX idx_location_updates_session_id ON location_updates(session_id);
CREATE INDEX idx_location_updates_received_at ON location_updates(received_at);
CREATE INDEX idx_location_updates_location ON location_updates(latitude, longitude);
CREATE INDEX idx_location_updates_vehicle_session ON location_updates(vehicle_id, session_id);