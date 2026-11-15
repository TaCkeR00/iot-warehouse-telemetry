CREATE TABLE IF NOT EXISTS warehouses (
  id INT AUTO_INCREMENT PRIMARY KEY,
  state VARCHAR(15),
  city VARCHAR(15),
  address VARCHAR(50)
);

CREATE TABLE IF NOT EXISTS sensors (
  id INT AUTO_INCREMENT PRIMARY KEY,
  model VARCHAR(50) NOT NULL,
  max_temperature INT NOT NULL,
  min_temperature INT NOT NULL,
  max_humidity INT NOT NULL,
  min_humidity INT NOT NULL,
  warehouse_id INT NOT NULL,
  FOREIGN KEY (warehouse_id) REFERENCES warehouses(id)
);

CREATE TABLE IF NOT EXISTS records (
  id INT AUTO_INCREMENT PRIMARY KEY,
  sensor_id INT NOT NULL,
  temperature INT NOT NULL,
  humidity INT NOT NULL,
  time DATETIME NOT NULL,
  FOREIGN KEY (sensor_id) REFERENCES sensors(id)
);

CREATE TABLE IF NOT EXISTS users (
  id INT AUTO_INCREMENT PRIMARY KEY,
  username VARCHAR(50) NOT NULL UNIQUE,
  password VARCHAR(255) NOT NULL,
  email VARCHAR(50) NOT NULL,
  is_admin BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE TABLE IF NOT EXISTS alerts (
  id INT AUTO_INCREMENT PRIMARY KEY,
  record_id INT NOT NULL,
  FOREIGN KEY (record_id) REFERENCES records(id)
);
