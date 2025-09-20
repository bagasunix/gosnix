-- Tabel tracking_sessions
CREATE TABLE tracking_sessions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    vehicle_id UUID NOT NULL,
    session_name VARCHAR(255) NOT NULL,
    start_time TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    end_time TIMESTAMP WITHOUT TIME ZONE,
    status VARCHAR(50) DEFAULT 'ACTIVE' NOT NULL,
    total_distance DECIMAL(10,3) DEFAULT 0.0,
    total_duration INTEGER DEFAULT 0,
    created_by INTEGER NOT NULL,
    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITHOUT TIME ZONE,
    deleted_at TIMESTAMP WITHOUT TIME ZONE,
    CONSTRAINT fk_tracking_sessions_vehicle FOREIGN KEY (vehicle_id) REFERENCES vehicles(id)
);

-- Index untuk tracking_sessions
CREATE INDEX idx_tracking_sessions_vehicle_id ON tracking_sessions(vehicle_id);
CREATE INDEX idx_tracking_sessions_created_by ON tracking_sessions(created_by);
CREATE INDEX idx_tracking_sessions_status ON tracking_sessions(status);
CREATE INDEX idx_tracking_sessions_start_time ON tracking_sessions(start_time);
CREATE INDEX idx_tracking_sessions_deleted_at ON tracking_sessions(deleted_at);
CREATE INDEX idx_tracking_sessions_time_range ON tracking_sessions(start_time, end_time);

ALTER TABLE tracking_sessions ADD CONSTRAINT fk_tracking_sessions_created_by FOREIGN KEY (created_by) REFERENCES customers(id);