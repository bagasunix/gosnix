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
    (99a6a003-7cac-43af-9f27-ca4c88c24027, 'Car', 'Kendaraan penumpang seperti sedan, hatchback, dan coupe'),
    (74e2f928-90a4-4658-80a9-3fc6d5459860, 'Motorcycle', 'Kendaraan bermotor beroda dua termasuk skuter dan motor sport'),
    (ef5a7207-54cf-4830-b8ec-652dbf3c1b91, 'Truck', 'Truk ringan, sedang, hingga berat untuk angkutan barang'),
    (9897f028-7ebc-46d7-9132-1d4ae724e3a8, 'Bus', 'Kendaraan untuk mengangkut banyak penumpang'),
    (9a1cf11b-b219-406e-a747-8347816d72d4, 'Van', 'Kendaraan serbaguna untuk angkut barang maupun penumpang'),
    (a6e65556-5429-4631-af30-2dc2f8968f91, 'SUV', 'Sport Utility Vehicle untuk keluarga dan medan off-road'),
    (358f4c5d-d3d6-46c6-b64d-3345747972b3, 'Pickup', 'Truk ringan dengan bak terbuka di bagian belakang'),
    (e7202487-acb0-4c38-bd05-e068c72ec70d, 'Heavy Equipment', 'Kendaraan konstruksi seperti ekskavator, buldoser, dan loader'),
    (e1a2df18-9f54-430a-9957-7bda92432442, 'Bicycle', 'Kendaraan beroda dua tanpa mesin (daya kayuh manusia)'),
    (516a6934-79ab-4a38-85a8-98d832efe295, 'Electric Vehicle', 'Kendaraan bertenaga baterai seperti mobil listrik dan sepeda listrik');
