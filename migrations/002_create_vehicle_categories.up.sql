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

-- Insert vehicle categories
INSERT INTO vehicle_categories (id, name, description)
VALUES
    (uuid_generate_v4(), 'Car', 'Kendaraan penumpang seperti sedan, hatchback, dan coupe'),
    (uuid_generate_v4(), 'Motorcycle', 'Kendaraan bermotor beroda dua termasuk skuter dan motor sport'),
    (uuid_generate_v4(), 'Truck', 'Truk ringan, sedang, hingga berat untuk angkutan barang'),
    (uuid_generate_v4(), 'Bus', 'Kendaraan untuk mengangkut banyak penumpang'),
    (uuid_generate_v4(), 'Van', 'Kendaraan serbaguna untuk angkut barang maupun penumpang'),
    (uuid_generate_v4(), 'SUV', 'Sport Utility Vehicle untuk keluarga dan medan off-road'),
    (uuid_generate_v4(), 'Pickup', 'Truk ringan dengan bak terbuka di bagian belakang'),
    (uuid_generate_v4(), 'Heavy Equipment', 'Kendaraan konstruksi seperti ekskavator, buldoser, dan loader'),
    (uuid_generate_v4(), 'Bicycle', 'Kendaraan beroda dua tanpa mesin (daya kayuh manusia)'),
    (uuid_generate_v4(), 'Electric Vehicle', 'Kendaraan bertenaga baterai seperti mobil listrik dan sepeda listrik');
