-- Create Users table
CREATE TABLE users (
user_id SERIAL PRIMARY KEY,
name TEXT NOT NULL
);

-- Seed Users table
INSERT INTO users (name) VALUES
('Alice Johnson'),
('Bob Smith'),
('Charlie Davis');

-- Create Rides table
CREATE TABLE rides (
ride_id SERIAL PRIMARY KEY,
source TEXT NOT NULL,
destination TEXT NOT NULL,
distance INT NOT NULL,
cost INT NOT NULL
);

-- Seed Rides table
INSERT INTO rides (source, destination, distance, cost) VALUES
('Downtown', 'Airport', 15, 150),
('City Center', 'Mall', 8, 80),
('Train Station', 'University', 12, 120);

-- Create Bookings table
CREATE TABLE bookings (
booking_id SERIAL PRIMARY KEY,
user_id INT REFERENCES users(user_id),
ride_id INT REFERENCES rides(ride_id),
time TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Seed Bookings table
INSERT INTO bookings (user_id, ride_id) VALUES
(1, 1),
(2, 2),
(3, 3);
